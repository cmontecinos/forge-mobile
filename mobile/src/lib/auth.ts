import Constants from "expo-constants";
import * as SecureStore from "expo-secure-store";

const API_URL =
  process.env.API_URL ||
  Constants.expoConfig?.extra?.apiUrl ||
  "http://localhost:8080";

// Types
export interface User {
  id: string;
  email: string;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}

export interface AuthError {
  error: string;
  message: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterCredentials {
  email: string;
  password: string;
}

// Token storage keys
const ACCESS_TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";
const USER_KEY = "user";

// Token storage helpers using SecureStore
export async function getAccessToken(): Promise<string | null> {
  return SecureStore.getItemAsync(ACCESS_TOKEN_KEY);
}

export async function getRefreshToken(): Promise<string | null> {
  return SecureStore.getItemAsync(REFRESH_TOKEN_KEY);
}

export async function getStoredUser(): Promise<User | null> {
  const user = await SecureStore.getItemAsync(USER_KEY);
  return user ? JSON.parse(user) : null;
}

export async function storeTokens(response: AuthResponse): Promise<void> {
  await SecureStore.setItemAsync(ACCESS_TOKEN_KEY, response.access_token);
  await SecureStore.setItemAsync(REFRESH_TOKEN_KEY, response.refresh_token);
  await SecureStore.setItemAsync(USER_KEY, JSON.stringify(response.user));
}

export async function clearTokens(): Promise<void> {
  await SecureStore.deleteItemAsync(ACCESS_TOKEN_KEY);
  await SecureStore.deleteItemAsync(REFRESH_TOKEN_KEY);
  await SecureStore.deleteItemAsync(USER_KEY);
}

// Auth API functions
export async function register(
  credentials: RegisterCredentials
): Promise<AuthResponse> {
  const response = await fetch(`${API_URL}/auth/register`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(credentials),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error((data as AuthError).message || "Registration failed");
  }

  await storeTokens(data as AuthResponse);
  return data as AuthResponse;
}

export async function login(
  credentials: LoginCredentials
): Promise<AuthResponse> {
  const response = await fetch(`${API_URL}/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(credentials),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error((data as AuthError).message || "Login failed");
  }

  await storeTokens(data as AuthResponse);
  return data as AuthResponse;
}

export async function logout(): Promise<void> {
  const token = await getAccessToken();

  if (token) {
    try {
      await fetch(`${API_URL}/auth/logout`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    } catch {
      // Best effort logout - ignore errors
    }
  }

  await clearTokens();
}

export async function refreshToken(): Promise<AuthResponse> {
  const refresh = await getRefreshToken();

  if (!refresh) {
    throw new Error("No refresh token available");
  }

  const response = await fetch(`${API_URL}/auth/refresh`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ refresh_token: refresh }),
  });

  const data = await response.json();

  if (!response.ok) {
    await clearTokens();
    throw new Error((data as AuthError).message || "Token refresh failed");
  }

  await storeTokens(data as AuthResponse);
  return data as AuthResponse;
}
