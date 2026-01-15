import { useEffect, useState } from "react";
import { View, Text, TouchableOpacity, ActivityIndicator } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useAuth } from "../hooks/useAuth";
import { checkHealth } from "../lib/api";

export function HomeScreen() {
  const { user, logout } = useAuth();
  const [healthStatus, setHealthStatus] = useState<string>("checking...");
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  useEffect(() => {
    checkHealth()
      .then((response) => setHealthStatus(response.status))
      .catch(() => setHealthStatus("offline"));
  }, []);

  const handleLogout = async () => {
    setIsLoggingOut(true);
    try {
      await logout();
    } finally {
      setIsLoggingOut(false);
    }
  };

  return (
    <SafeAreaView className="flex-1 bg-white">
      <View className="flex-1 items-center justify-center px-6">
        <Text className="text-3xl font-bold text-gray-900 mb-2">
          Welcome to {"{{.ProjectName}}"}
        </Text>
        <Text className="text-lg text-gray-600 mb-8">
          Expo + Go + Supabase
        </Text>

        {/* User Info */}
        <View className="bg-white border border-gray-200 rounded-lg px-6 py-4 mb-6 w-full max-w-sm">
          <Text className="text-sm text-gray-500 text-center mb-1">
            Logged in as
          </Text>
          <Text className="text-base font-medium text-gray-900 text-center mb-4">
            {user?.email}
          </Text>
          <TouchableOpacity
            className={`py-2 px-4 rounded-lg ${
              isLoggingOut ? "bg-red-400" : "bg-red-600"
            }`}
            onPress={handleLogout}
            disabled={isLoggingOut}
          >
            {isLoggingOut ? (
              <ActivityIndicator color="white" size="small" />
            ) : (
              <Text className="text-white text-center font-medium">Logout</Text>
            )}
          </TouchableOpacity>
        </View>

        {/* Health Check */}
        <View className="bg-gray-100 rounded-lg px-4 py-3">
          <Text className="text-sm text-gray-500">
            Backend:{" "}
            <Text
              className={`font-medium ${
                healthStatus === "healthy" ? "text-green-600" : "text-red-600"
              }`}
            >
              {healthStatus}
            </Text>
          </Text>
        </View>
      </View>
    </SafeAreaView>
  );
}
