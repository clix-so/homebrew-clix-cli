import { Stack } from 'expo-router';
import { useEffect } from 'react';

export default function RootLayout() {
  return (
    <Stack>
      <Stack.Screen name="index" />
    </Stack>
  );
}