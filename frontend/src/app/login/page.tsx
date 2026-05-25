'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useLogin } from '@/hooks/useAuth';
import Link from 'next/link';

export default function LoginPage() {
  const router = useRouter();
  const loginMutation = useLogin();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMessage('');

    try {
      const response = await loginMutation.mutateAsync({ email, password });
      localStorage.setItem('access_token', response.access_token);
      
      // Dispatch event to update Navbar immediately
      window.dispatchEvent(new Event('auth-change'));
      
      router.push('/feed');
    } catch (error: any) {
      const errorStr = error.message || String(error);
      if (errorStr.includes('401') || errorStr.toLowerCase().includes('invalid credentials')) {
        setErrorMessage('E-posta veya şifre hatalı. Lütfen bilgilerinizi kontrol edin.');
      } else if (errorStr.includes('400')) {
        setErrorMessage('Geçersiz bir giriş yaptınız. Lütfen e-posta ve şifre alanlarını kontrol edin.');
      } else {
        setErrorMessage('Giriş yapılırken beklenmedik bir hata oluştu. Lütfen daha sonra tekrar deneyin.');
      }
    }
  };

  return (
    <div className="max-w-sm mx-auto mt-12 md:mt-24 p-6 bg-white border border-gray-200 rounded-xl shadow-sm">
      <div className="mb-8 text-center">
        <h1 className="text-2xl font-bold text-gray-900">Giriş Yap</h1>
        <p className="text-sm text-gray-500 mt-2">Hesabınıza erişmek için bilgilerinizi girin</p>
      </div>
      
      {errorMessage && (
        <div className="mb-6 p-3 bg-red-50 text-red-700 text-sm rounded-md border border-red-100">
          {errorMessage}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-5">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1.5" htmlFor="email">
            E-posta adresi
          </label>
          <input 
            id="email"
            type="email" 
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            autoComplete="email"
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black focus:border-transparent text-gray-900 transition-shadow"
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1.5" htmlFor="password">
            Şifre
          </label>
          <input 
            id="password"
            type="password" 
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            autoComplete="current-password"
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black focus:border-transparent text-gray-900 transition-shadow"
          />
        </div>

        <button 
          type="submit" 
          disabled={loginMutation.isPending}
          className="w-full py-2.5 px-4 bg-black text-white font-medium rounded-lg hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-black disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {loginMutation.isPending ? 'Giriş yapılıyor...' : 'Giriş Yap'}
        </button>
      </form>

      <div className="mt-6 text-center text-sm text-gray-600">
        Hesabınız yok mu?{' '}
        <Link href="/register" className="font-medium text-black hover:underline">
          Kayıt ol
        </Link>
      </div>
    </div>
  );
}
