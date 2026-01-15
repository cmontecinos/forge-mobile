import { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
} from "react-native";
import { useAuth } from "../hooks/useAuth";
import type { NativeStackNavigationProp } from "@react-navigation/native-stack";

type LoginScreenProps = {
  navigation: NativeStackNavigationProp<any>;
};

export function LoginScreen({ navigation }: LoginScreenProps) {
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async () => {
    setError("");

    if (!email || !password) {
      setError("Email and password are required");
      return;
    }

    setIsLoading(true);

    try {
      await login({ email, password });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === "ios" ? "padding" : "height"}
      className="flex-1 bg-gray-50"
    >
      <View className="flex-1 justify-center px-6">
        <Text className="text-2xl font-bold text-center text-gray-900 mb-8">
          Sign In
        </Text>

        {error ? (
          <View className="bg-red-50 border border-red-200 rounded-lg px-4 py-3 mb-4">
            <Text className="text-red-600 text-center">{error}</Text>
          </View>
        ) : null}

        <View className="mb-4">
          <Text className="text-sm font-medium text-gray-700 mb-1">Email</Text>
          <TextInput
            className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-white text-gray-900"
            placeholder="you@example.com"
            value={email}
            onChangeText={setEmail}
            autoCapitalize="none"
            keyboardType="email-address"
            editable={!isLoading}
          />
        </View>

        <View className="mb-6">
          <Text className="text-sm font-medium text-gray-700 mb-1">
            Password
          </Text>
          <TextInput
            className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-white text-gray-900"
            placeholder="Enter your password"
            value={password}
            onChangeText={setPassword}
            secureTextEntry
            editable={!isLoading}
          />
        </View>

        <TouchableOpacity
          className={`w-full py-3 rounded-lg ${
            isLoading ? "bg-blue-400" : "bg-blue-600"
          }`}
          onPress={handleSubmit}
          disabled={isLoading}
        >
          {isLoading ? (
            <ActivityIndicator color="white" />
          ) : (
            <Text className="text-white text-center font-semibold">
              Sign In
            </Text>
          )}
        </TouchableOpacity>

        <TouchableOpacity
          className="mt-4"
          onPress={() => navigation.navigate("Register")}
        >
          <Text className="text-center text-gray-600">
            Don't have an account?{" "}
            <Text className="text-blue-600">Sign up</Text>
          </Text>
        </TouchableOpacity>
      </View>
    </KeyboardAvoidingView>
  );
}
