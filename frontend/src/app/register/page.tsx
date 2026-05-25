'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useRegister } from '@/hooks/useAuth';
import Link from 'next/link';

export default function RegisterPage() {
  const router = useRouter();
  const registerMutation = useRegister();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [displayName, setDisplayName] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMessage('');

    try {
      const response = await registerMutation.mutateAsync({ 
        email, 
        password, 
        display_name: displayName 
      });
      localStorage.setItem('access_token', response.access_token);
      
      // Dispatch event to update Navbar immediately
      window.dispatchEvent(new Event('auth-change'));
      
      router.push('/feed');
    } catch (error: any) {
      const errorStr = error.message || String(error);
      if (errorStr.includes('409') || errorStr.toLowerCase().includes('already registered')) {
        setErrorMessage('Bu e-posta adresi sistemde zaten kayıtlı. Lütfen giriş yapmayı deneyin.');
      } else if (errorStr.includes('400')) {
        setErrorMessage('Girdiğiniz bilgiler geçersiz. Lütfen alanları eksiksiz ve doğru doldurun.');
      } else {
        setErrorMessage('Kayıt olurken beklenmedik bir hata oluştu. Lütfen daha sonra tekrar deneyin.');
      }
    }
  };

  return (
    <div className="max-w-sm mx-auto mt-12 md:mt-24 p-6 bg-white border border-gray-200 rounded-xl shadow-sm">
      <div className="mb-8 text-center">
        <h1 className="text-2xl font-bold text-gray-900">Kayıt Ol</h1>
        <p className="text-sm text-gray-500 mt-2">PulsePoll'a katılmak için bir hesap oluşturun</p>
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
          <label className="block text-sm font-medium text-gray-700 mb-1.5" htmlFor="displayName">
            Görüntülenme Adı
          </label>
          <input 
            id="displayName"
            type="text" 
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
            required
            autoComplete="name"
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
            autoComplete="new-password"
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black focus:border-transparent text-gray-900 transition-shadow"
          />
        </div>

        <button 
          type="submit" 
          disabled={registerMutation.isPending}
          className="w-full py-2.5 px-4 bg-black text-white font-medium rounded-lg hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-black disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {registerMutation.isPending ? 'Hesap oluşturuluyor...' : 'Kayıt Ol'}
        </button>
      </form>

      <div className="mt-6 text-center text-sm text-gray-600">
        Zaten bir hesabınız var mı?{' '}
        <Link href="/login" className="font-medium text-black hover:underline">
          Giriş yap
        </Link>
      </div>
    </div>
  );
}
