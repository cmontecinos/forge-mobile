import "./global.css";
import { StatusBar } from "expo-status-bar";
import { useEffect, useState } from "react";
import { Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { checkHealth } from "./src/lib/api";

export default function App() {
  const [healthStatus, setHealthStatus] = useState<string>("checking...");

  useEffect(() => {
    checkHealth()
      .then((response) => setHealthStatus(response.status))
      .catch(() => setHealthStatus("offline"));
  }, []);

  return (
    <SafeAreaView className="flex-1 bg-white">
      <View className="flex-1 items-center justify-center px-4">
        <Text className="text-3xl font-bold text-gray-900 mb-2">
          Welcome to {"{{.ProjectName}}"}
        </Text>
        <Text className="text-lg text-gray-600 mb-8">
          Expo + Go + Supabase
        </Text>
        <View className="bg-gray-100 rounded-lg px-4 py-3">
          <Text className="text-sm text-gray-500">
            Backend: <Text className="font-medium text-gray-700">{healthStatus}</Text>
          </Text>
        </View>
      </View>
      <StatusBar style="auto" />
    </SafeAreaView>
  );
}
