import { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
} from "react-native";
import { useAuth } from "../hooks/useAuth";
import type { NativeStackNavigationProp } from "@react-navigation/native-stack";

type RegisterScreenProps = {
  navigation: NativeStackNavigationProp<any>;
};

export function RegisterScreen({ navigation }: RegisterScreenProps) {
  const { register } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async () => {
    setError("");

    if (!email || !password || !confirmPassword) {
      setError("All fields are required");
      return;
    }

    if (password.length < 6) {
      setError("Password must be at least 6 characters");
      return;
    }

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    setIsLoading(true);

    try {
      await register({ email, password });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === "ios" ? "padding" : "height"}
      className="flex-1 bg-gray-50"
    >
      <ScrollView
        contentContainerStyle={{ flexGrow: 1 }}
        keyboardShouldPersistTaps="handled"
      >
        <View className="flex-1 justify-center px-6 py-8">
          <Text className="text-2xl font-bold text-center text-gray-900 mb-8">
            Create Account
          </Text>

          {error ? (
            <View className="bg-red-50 border border-red-200 rounded-lg px-4 py-3 mb-4">
              <Text className="text-red-600 text-center">{error}</Text>
            </View>
          ) : null}

          <View className="mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-1">
              Email
            </Text>
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

          <View className="mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-1">
              Password
            </Text>
            <TextInput
              className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-white text-gray-900"
              placeholder="At least 6 characters"
              value={password}
              onChangeText={setPassword}
              secureTextEntry
              editable={!isLoading}
            />
          </View>

          <View className="mb-6">
            <Text className="text-sm font-medium text-gray-700 mb-1">
              Confirm Password
            </Text>
            <TextInput
              className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-white text-gray-900"
              placeholder="Confirm your password"
              value={confirmPassword}
              onChangeText={setConfirmPassword}
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
                Create Account
              </Text>
            )}
          </TouchableOpacity>

          <TouchableOpacity
            className="mt-4"
            onPress={() => navigation.navigate("Login")}
          >
            <Text className="text-center text-gray-600">
              Already have an account?{" "}
              <Text className="text-blue-600">Sign in</Text>
            </Text>
          </TouchableOpacity>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}
