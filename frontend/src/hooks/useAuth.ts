import { useMutation, useQuery } from '@tanstack/react-query';

import { getMe, login, logout, register } from '@/services/auth';

type LoginVariables = {
  email: string;
  password: string;
};

type RegisterVariables = {
  email: string;
  password: string;
  display_name: string;
};

export function useMe(access_token: string) {
  return useQuery({
    queryKey: ['me'],
    queryFn: () => getMe(access_token),
    enabled: access_token.length > 0,
  });
}

export function useLogin() {
  return useMutation({
    mutationFn: ({ email, password }: LoginVariables) => login(email, password),
  });
}

export function useRegister() {
  return useMutation({
    mutationFn: ({ email, password, display_name }: RegisterVariables) =>
      register(email, password, display_name),
  });
}

export function useLogout() {
  return useMutation({
    mutationFn: () => logout(),
  });
}
