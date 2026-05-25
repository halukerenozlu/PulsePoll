'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { useLogout } from '@/hooks/useAuth';

export function Navbar() {
  const router = useRouter();
  const [hasToken, setHasToken] = useState(false);
  const logoutMutation = useLogout();

  useEffect(() => {
    const checkToken = () => {
      setHasToken(!!localStorage.getItem('access_token'));
    };
    
    // Check token on mount
    checkToken();
    
    // Listen for custom auth events to update UI dynamically without reload
    window.addEventListener('auth-change', checkToken);
    
    return () => {
      window.removeEventListener('auth-change', checkToken);
    };
  }, []);

  const handleLogout = async () => {
    try {
      await logoutMutation.mutateAsync();
    } catch (e) {
      console.error('Logout failed on server, continuing local logout', e);
    } finally {
      localStorage.removeItem('access_token');
      window.dispatchEvent(new Event('auth-change'));
      router.push('/');
    }
  };

  return (
    <nav className="flex items-center justify-between p-4 border-b border-gray-200 bg-white sticky top-0 z-10">
      <Link href="/" className="text-xl font-bold text-black tracking-tight">
        PulsePoll
      </Link>
      <div className="flex gap-3 md:gap-4 items-center">
        {hasToken ? (
          <>
            <Link
              href="/surveys/create"
              className="text-sm font-medium px-4 py-2 bg-black text-white hover:bg-gray-800 rounded-md transition-colors"
            >
              Anket Oluştur
            </Link>
            <button
              onClick={handleLogout}
              disabled={logoutMutation.isPending}
              className="text-sm font-medium px-4 py-2 bg-gray-100 text-gray-900 hover:bg-gray-200 rounded-md transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {logoutMutation.isPending ? 'Çıkılıyor...' : 'Çıkış Yap'}
            </button>
          </>
        ) : (
          <>
            <Link
              href="/login"
              className="text-sm font-medium px-4 py-2 text-gray-700 hover:text-black transition-colors"
            >
              Giriş Yap
            </Link>
            <Link
              href="/register"
              className="text-sm font-medium px-4 py-2 bg-black text-white hover:bg-gray-800 rounded-md transition-colors"
            >
              Kayıt Ol
            </Link>
          </>
        )}
      </div>
    </nav>
  );
}
