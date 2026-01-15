import "./global.css";
import { StatusBar } from "expo-status-bar";
import { NavigationContainer } from "@react-navigation/native";
import { AuthProvider } from "./src/contexts/AuthContext";
import { AuthNavigator } from "./src/navigation/AuthNavigator";

export default function App() {
  return (
    <AuthProvider>
      <NavigationContainer>
        <AuthNavigator />
        <StatusBar style="auto" />
      </NavigationContainer>
    </AuthProvider>
  );
}
