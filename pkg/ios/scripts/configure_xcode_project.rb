#!/usr/bin/env ruby

require 'xcodeproj'
require 'fileutils'
require 'optparse'
require 'json'

# Parse command-line options
options = {}
OptionParser.new do |opts|
  opts.banner = "Usage: configure_xcode_project.rb [options]"

  opts.on("--project-path PATH", "Path to the .xcodeproj file") do |path|
    options[:project_path] = path
  end

  opts.on("--app-group-id ID", "App Group ID (e.g., group.clix.myproject)") do |id|
    options[:app_group_id] = id
  end

  opts.on("--main-target NAME", "Main app target name") do |name|
    options[:main_target] = name
  end

  opts.on("--extension-target NAME", "NotificationServiceExtension target name (defaults to 'NotificationServiceExtension')") do |name|
    options[:extension_target] = name
  end

  opts.on("--framework-name NAME", "Framework to add to targets (defaults to 'Clix')") do |name|
    options[:framework_name] = name
  end

  opts.on("--verbose", "Enable verbose output") do
    options[:verbose] = true
  end
end.parse!

# Set defaults for optional parameters
options[:extension_target] ||= "NotificationServiceExtension"
options[:framework_name] ||= "Clix"

# Validate required options
required = [:project_path, :app_group_id]
missing = required.select { |param| options[param].nil? }
if !missing.empty?
  $stderr.puts "Missing required parameters: #{missing.join(', ')}"
  exit 1
end

# Define log function - ALWAYS use stderr for logs
def log(message, options)
  # IMPORTANT: Use stderr for ALL logging to avoid interfering with JSON output
  $stderr.puts message if options[:verbose]
end

# Output final result as JSON to stdout
def result(success, message, data = {})
  # CRITICAL: This is the ONLY output to stdout
  # All other output must use stderr to avoid breaking JSON parsing
  $stdout.puts JSON.generate({success: success, message: message, data: data})
end

begin
  # IMPORTANT: All output to stdout must be valid JSON only
  # All log messages must go to stderr
  log("Opening Xcode project: #{options[:project_path]}", options)
  project = Xcodeproj::Project.open(options[:project_path])
  
  # Find targets
  main_target = nil
  extension_target = nil
  
  if options[:main_target]
    main_target = project.targets.find { |t| t.name == options[:main_target] }
    if main_target.nil?
      result(false, "Main target '#{options[:main_target]}' not found")
      exit 1
    end
  else
    # Try to find the main app target (usually the first application target)
    main_target = project.targets.find { |t| t.product_type == "com.apple.product-type.application" }
    if main_target.nil?
      result(false, "Could not automatically determine main app target")
      exit 1
    end
  end
  
  extension_target = project.targets.find { |t| t.name == options[:extension_target] }
  if extension_target.nil?
    result(false, "Extension target '#{options[:extension_target]}' not found")
    exit 1
  end
  
  log("Found main target: #{main_target.name}", options)
  log("Found extension target: #{extension_target.name}", options)
  
  # Get project directory
  project_dir = File.dirname(options[:project_path])
  
  # Track modifications for reporting
  main_target_modified = false
  extension_modified = false
  main_framework_added = false
  extension_framework_added = false
  
  #-----------------------------------------
  # 1. Configure main app target entitlements
  #-----------------------------------------
  main_target.build_configurations.each do |config|
    # Try to find the main target's directory
    main_target_dir = nil
    
    # Try standard convention: project_dir/TargetName
    standard_dir = File.join(project_dir, main_target.name)
    if Dir.exist?(standard_dir)
      main_target_dir = standard_dir
    else
      # Try to determine from build files
      build_file_dirs = {}
      
      # Check source build phase files
      main_target.source_build_phase.files_references.each do |file_ref|
        if file_ref.real_path && File.exist?(file_ref.real_path)
          dir = File.dirname(file_ref.real_path.to_s)
          build_file_dirs[dir] ||= 0
          build_file_dirs[dir] += 1
        end
      end
      
      # Check resource build phase files as well
      main_target.resources_build_phase.files_references.each do |file_ref|
        if file_ref.real_path && File.exist?(file_ref.real_path)
          dir = File.dirname(file_ref.real_path.to_s)
          build_file_dirs[dir] ||= 0
          build_file_dirs[dir] += 1
        end
      end
      
      # Find the most common directory if we have build files
      if !build_file_dirs.empty?
        main_target_dir = build_file_dirs.max_by { |dir, count| count }[0]
      else
        # Fallback to project root
        main_target_dir = project_dir
      end
    end
    
    # Determine entitlements file path
    entitlements_path = nil
    if config.build_settings['CODE_SIGN_ENTITLEMENTS']
      # Get path from existing setting
      relative_path = config.build_settings['CODE_SIGN_ENTITLEMENTS']
      full_path = File.join(project_dir, relative_path)
      
      if File.exist?(full_path)
        # Use existing path
        entitlements_path = full_path
        log("Using existing main target entitlements file at: #{entitlements_path}", options)
      else
        # Create in target directory instead
        entitlements_path = File.join(main_target_dir, File.basename(relative_path))
      end
    else
      # Create a new entitlements file in the target directory
      target_subdir = File.join(project_dir, main_target.name)
      
      # Ensure the target directory exists
      unless Dir.exist?(target_subdir)
        FileUtils.mkdir_p(target_subdir)
        log("Created target directory: #{target_subdir}", options)
      end
      
      entitlements_filename = "#{main_target.name}.entitlements"
      entitlements_path = File.join(target_subdir, entitlements_filename)
      
      # Set path relative to the Xcode project directory
      rel_path = Pathname.new(entitlements_path).relative_path_from(Pathname.new(project_dir)).to_s
      config.build_settings['CODE_SIGN_ENTITLEMENTS'] = rel_path
      log("Created new entitlements file at: #{entitlements_path}", options)
      main_target_modified = true
    end
    
    # Create or update entitlements file
    entitlements = File.exist?(entitlements_path) ? 
      Xcodeproj::Plist.read_from_path(entitlements_path) : {}
    
    # Add App Groups entitlement
    if !entitlements['com.apple.security.application-groups']
      entitlements['com.apple.security.application-groups'] = []
      main_target_modified = true
    end
    
    # Add specific app group if not already present
    if !entitlements['com.apple.security.application-groups'].include?(options[:app_group_id])
      entitlements['com.apple.security.application-groups'] << options[:app_group_id]
      main_target_modified = true
    end
    
    # Add Push Notifications entitlement
    if !entitlements['aps-environment']
      entitlements['aps-environment'] = 'development'
      main_target_modified = true
    end
    
    # Write entitlements back to file
    Xcodeproj::Plist.write_to_path(entitlements, entitlements_path)
  end
  
  #-----------------------------------------
  # 2. Configure extension target entitlements
  #-----------------------------------------
  extension_target.build_configurations.each do |config|
    # Try to find the extension target's directory
    extension_target_dir = nil
    
    # Try standard convention: project_dir/TargetName
    standard_dir = File.join(project_dir, extension_target.name)
    if Dir.exist?(standard_dir)
      extension_target_dir = standard_dir
    else
      # Try to create the directory if it doesn't exist
      begin
        FileUtils.mkdir_p(standard_dir)
        extension_target_dir = standard_dir
      rescue => e
        log("Failed to create extension target directory: #{e.message}", options)
        
        # Try to determine from build files
        build_file_dirs = {}
        
        # Check source and resource build phases
        [extension_target.source_build_phase, extension_target.resources_build_phase].each do |phase|
          phase.files_references.each do |file_ref|
            if file_ref.real_path && File.exist?(file_ref.real_path)
              dir = File.dirname(file_ref.real_path.to_s)
              build_file_dirs[dir] ||= 0
              build_file_dirs[dir] += 1
            end
          end
        end
        
        # Find most common directory or use project root
        extension_target_dir = if !build_file_dirs.empty?
                                build_file_dirs.max_by { |dir, count| count }[0]
                              else
                                project_dir
                              end
      end
    end
    
    # Determine entitlements file path
    entitlements_path = nil
    if config.build_settings['CODE_SIGN_ENTITLEMENTS']
      # Get path from existing setting
      relative_path = config.build_settings['CODE_SIGN_ENTITLEMENTS']
      full_path = File.join(project_dir, relative_path)
      
      if File.exist?(full_path)
        entitlements_path = full_path
        log("Using existing extension entitlements file at: #{entitlements_path}", options)
      else
        entitlements_path = File.join(extension_target_dir, File.basename(relative_path))
      end
    else
      # Create new entitlements file
      target_subdir = File.join(project_dir, extension_target.name)
      
      unless Dir.exist?(target_subdir)
        FileUtils.mkdir_p(target_subdir)
      end
      
      entitlements_filename = "#{extension_target.name}.entitlements"
      entitlements_path = File.join(target_subdir, entitlements_filename)
      
      rel_path = Pathname.new(entitlements_path).relative_path_from(Pathname.new(project_dir)).to_s
      config.build_settings['CODE_SIGN_ENTITLEMENTS'] = rel_path
      log("Created new extension entitlements file at: #{entitlements_path}", options)
      extension_modified = true
    end
    
    # Create or update entitlements file
    entitlements = File.exist?(entitlements_path) ? 
      Xcodeproj::Plist.read_from_path(entitlements_path) : {}
    
    # Add App Groups entitlement
    if !entitlements['com.apple.security.application-groups']
      entitlements['com.apple.security.application-groups'] = []
      extension_modified = true
    end
    
    # Add specific app group if not already present
    if !entitlements['com.apple.security.application-groups'].include?(options[:app_group_id])
      entitlements['com.apple.security.application-groups'] << options[:app_group_id]
      extension_modified = true
    end
    
    # Write entitlements back to file
    Xcodeproj::Plist.write_to_path(entitlements, entitlements_path)
  end

  #-----------------------------------------
  # 3. Add Clix framework to main app target
  #-----------------------------------------
  framework_name = options[:framework_name]
  
  # Check if framework is already added to main target
  if !main_target.frameworks_build_phase.files.find { |f| f.file_ref && f.file_ref.respond_to?(:display_name) && f.file_ref.display_name == framework_name }
    # Look for framework reference
    framework_ref = nil
    
    # Check in frameworks group first
    if project.frameworks_group
      project.frameworks_group.children.each do |child|
        if child.display_name == framework_name
          framework_ref = child
          log("Found #{framework_name} in frameworks group", options)
          break
        end
      end
    end
    
    # If not found, we might be using SPM or CocoaPods, in which case we need to find it in the project
    if framework_ref.nil?
      # Try to find in any group (including Pods)
      pods_framework_found = false
      project.groups.each do |group|
        # Check if this is a Pods group (CocoaPods)
        if group.display_name == "Pods" || group.path =~ /Pods/
          log("Checking CocoaPods group: #{group.display_name}", options)
          # Look for the framework in Pods subgroups
          group.recursive_children.each do |child|
            if child.display_name == framework_name && (child.path.end_with?(".framework") || child.path =~ /#{framework_name}/)
              framework_ref = child
              pods_framework_found = true
              log("Found #{framework_name} framework in CocoaPods", options)
              break
            end
          end
        else
          # Regular group search
          group.recursive_children.each do |child|
            if child.display_name == framework_name && child.path.end_with?(".framework")
              framework_ref = child
              log("Found #{framework_name} framework in group: #{group.display_name}", options)
              break
            end
          end
        end
        break if framework_ref
      end
      
      # If found in CocoaPods, add it to the frameworks build phase
      if pods_framework_found && framework_ref
        extension_target.frameworks_build_phase.add_file_reference(framework_ref)
        extension_framework_added = true
        log("Added CocoaPods framework '#{framework_name}' to extension target", options)
      end
    end
    
    # If still not found, look for it as SPM dependency
    if framework_ref.nil? && !extension_framework_added
      # First check if the framework is already added as a dependency
      spm_dep_found = false
      
      # Check target dependencies for the framework
      extension_target.dependencies.each do |dep|
        if (dep.respond_to?(:target) && dep.target && dep.target.name == framework_name) ||
           (dep.respond_to?(:target_proxy) && dep.target_proxy && 
            dep.target_proxy.respond_to?(:remote_info) && dep.target_proxy.remote_info == framework_name)
          spm_dep_found = true
          log("#{framework_name} already added as dependency to extension target", options)
          break
        end
      end
      
      if spm_dep_found
        extension_framework_added = true
      else
        # Try to find SPM package references in the project
        spm_package = nil
        spm_product = nil
        
        # Check if project has any SPM package references
        if project.respond_to?(:package_references) && !project.package_references.empty?
          log("Project has #{project.package_references.count} SPM package references", options)
          
          # Look for a package that might contain our framework
          project.package_references.each do |pkg_ref|
            # Check if this package has our framework as a product
            if pkg_ref.respond_to?(:package_products)
              pkg_ref.package_products.each do |product|
                if product.respond_to?(:product_name) && product.product_name == framework_name
                  spm_package = pkg_ref
                  spm_product = product
                  log("Found #{framework_name} in SPM package: #{pkg_ref.name || pkg_ref.url}", options)
                  break
                end
              end
            end
            break if spm_package
          end
          
          # If we found the package and product, add it as a dependency
          if spm_package && spm_product
            # Create a dependency on the SPM product - must be a proper object, not just a string
            # Add product dependency correctly, ensuring we have the proper PBXTargetDependency object
            dependency = nil
            project.targets.each do |target|
              if target.respond_to?(:product_name) && target.product_name == framework_name
                dependency = extension_target.add_dependency(target)
                break
              end
            end
            
            if dependency
              extension_framework_added = true
              log("Added SPM dependency on '#{framework_name}' for extension target", options)
            else
              log("Found SPM product for '#{framework_name}' but could not create dependency", options)
            end
          end
        else
          log("No SPM packages found in project", options)
        end
        
        # If SPM approach didn't work, try to find a target with the same name
        if !extension_framework_added
          clix_target = project.targets.find { |t| t.name == framework_name }
          if clix_target
            # Use the target object, not just a string name
            extension_target.add_dependency(clix_target)
            extension_framework_added = true
            log("Added target dependency on '#{framework_name}' for extension target", options)
          else
            log("[WARN] Could not find a target, SPM package, or CocoaPods framework named '#{framework_name}' to add as dependency to extension target. Please add it manually in Xcode if needed.", options)
          end
        end
      end
    end
    
    # If we found a framework reference but haven't added it yet, add it to the extension target
    if framework_ref && !extension_framework_added
      extension_target.frameworks_build_phase.add_file_reference(framework_ref)
      extension_framework_added = true
      log("Added framework reference '#{framework_name}' to extension target", options)
    end
  else
    log("#{framework_name} framework already added to extension target", options)
    extension_framework_added = true
  end
  
  # Save the project
  project.save
  
  # Return success status and modifications made
  modifications = {
    main_target_modified: main_target_modified,
    extension_modified: extension_modified,
    main_framework_added: main_framework_added,
    extension_framework_added: extension_framework_added
  }
  
  result(true, "Xcode project configuration completed successfully", modifications)
  
rescue => e
  result(false, "Error: #{e.message}\n#{e.backtrace.join("\n")}")
end
