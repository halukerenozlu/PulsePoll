import { api } from './api';
import type { AuthUser, AuthResponse, LogoutResponse } from '@/types';

export function register(
  email: string,
  password: string,
  display_name: string,
): Promise<AuthResponse> {
  return api<AuthResponse>('/api/v1/auth/register', {
    method: 'POST',
    body: JSON.stringify({ email, password, display_name }),
    credentials: 'include',
  });
}

export function login(email: string, password: string): Promise<AuthResponse> {
  return api<AuthResponse>('/api/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
    credentials: 'include',
  });
}

export function logout(): Promise<LogoutResponse> {
  return api<LogoutResponse>('/api/v1/auth/logout', {
    method: 'POST',
    credentials: 'include',
  });
}

export function refreshToken(): Promise<AuthResponse> {
  return api<AuthResponse>('/api/v1/auth/refresh', {
    method: 'POST',
    credentials: 'include',
  });
}

export function getMe(access_token: string): Promise<AuthUser> {
  return api<AuthUser>('/api/v1/me', {
    headers: { Authorization: `Bearer ${access_token}` },
  });
}
