import "./global.css";
import { StatusBar } from "expo-status-bar";
import { useEffect, useState } from "react";
import { SafeAreaView, Text, View } from "react-native";
import { checkHealth } from "./src/lib/api";

export default function App() {
  const [healthStatus, setHealthStatus] = useState<string>("checking...");

  useEffect(() => {
    checkHealth()
      .then((data) => setHealthStatus(data.status))
      .catch((err) => setHealthStatus(`error: ${err.message}`));
  }, []);

  return (
    <SafeAreaView className="flex-1 bg-white">
      <View className="flex-1 items-center justify-center px-4">
        <Text className="text-2xl font-bold text-gray-900 mb-4">
          Welcome to {"{{.ProjectName}}"}
        </Text>
        <View className="bg-gray-100 rounded-lg p-4 w-full max-w-sm">
          <Text className="text-sm text-gray-600 text-center">
            Backend Status
          </Text>
          <Text className="text-lg font-medium text-center mt-1">
            {healthStatus}
          </Text>
        </View>
      </View>
      <StatusBar style="auto" />
    </SafeAreaView>
  );
}
