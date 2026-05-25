'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useFeed } from '@/hooks/useSurveys';
import type { SurveyPhase } from '@/types';

function getTimeRemainingText(voteEndsAt: string) {
  const diffMs = new Date(voteEndsAt).getTime() - Date.now();
  if (diffMs <= 0) return 'Süre doldu';
  
  const diffMins = Math.floor(diffMs / (1000 * 60));
  if (diffMins < 60) return `${diffMins} dakika kaldı`;
  
  const diffHours = Math.floor(diffMins / 60);
  if (diffHours < 24) return `${diffHours} saat kaldı`;
  
  const diffDays = Math.floor(diffHours / 24);
  return `${diffDays} gün kaldı`;
}

function PhaseBadge({ phase }: { phase: SurveyPhase }) {
  if (phase === 'VOTING') {
    return (
      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
        Oylama Devam Ediyor
      </span>
    );
  }
  if (phase === 'RESULTS') {
    return (
      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
        Sonuçlar
      </span>
    );
  }
  return (
    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
      Süresi Doldu
    </span>
  );
}

export default function FeedPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedSearch(searchTerm);
    }, 300);
    return () => clearTimeout(handler);
  }, [searchTerm]);

  const { data, isLoading, isError } = useFeed(
    debouncedSearch ? { search: debouncedSearch } : undefined
  );

  return (
    <div className="py-6">
      <div className="mb-8 flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Anket Akışı</h1>
          <p className="text-sm text-gray-500 mt-1">En son eklenen açık anketleri keşfedin</p>
        </div>
        <div className="w-full md:w-72">
          <div className="relative">
            <input
              type="text"
              placeholder="Anketlerde ara..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-4 pr-10 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black focus:border-transparent text-sm text-gray-900"
            />
            {searchTerm && (
              <button 
                onClick={() => setSearchTerm('')}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 focus:outline-none"
                aria-label="Aramayı temizle"
              >
                &times;
              </button>
            )}
          </div>
        </div>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {[1, 2, 3, 4].map(i => (
            <div key={i} className="p-5 border border-gray-200 rounded-xl animate-pulse bg-white">
              <div className="h-5 bg-gray-200 rounded w-3/4 mb-3"></div>
              <div className="h-4 bg-gray-200 rounded w-full mb-2"></div>
              <div className="h-4 bg-gray-200 rounded w-5/6 mb-4"></div>
              <div className="flex gap-2 mt-4">
                <div className="h-5 bg-gray-200 rounded w-20"></div>
                <div className="h-5 bg-gray-200 rounded w-24"></div>
              </div>
            </div>
          ))}
        </div>
      ) : isError ? (
        <div className="p-8 text-center bg-red-50 border border-red-100 rounded-xl">
          <p className="text-red-700">Anketler yüklenemedi, tekrar deneyin.</p>
        </div>
      ) : !data?.items || data.items.length === 0 ? (
        <div className="p-12 text-center border border-gray-200 border-dashed rounded-xl bg-white">
          <p className="text-gray-500 text-lg">Henüz anket yok</p>
          {debouncedSearch && (
            <p className="text-gray-400 text-sm mt-2">Arama kriterlerinize uyan bir anket bulunamadı.</p>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {data.items.map((survey) => (
            <Link 
              key={survey.id} 
              href={`/surveys/${survey.id}`}
              className="block p-5 border border-gray-200 rounded-xl hover:border-black hover:shadow-sm transition-all bg-white group"
            >
              <h2 className="text-lg font-semibold text-gray-900 line-clamp-2 group-hover:text-black">
                {survey.title}
              </h2>
              {survey.description && (
                <p className="text-sm text-gray-600 mt-2 line-clamp-2">
                  {survey.description}
                </p>
              )}
              
              <div className="mt-4 flex flex-wrap items-center gap-2">
                <PhaseBadge phase={survey.phase} />
                {survey.phase === 'VOTING' && survey.vote_ends_at && (
                  <span className="text-xs text-gray-500 font-medium bg-gray-100 px-2 py-0.5 rounded-md">
                    {getTimeRemainingText(survey.vote_ends_at)}
                  </span>
                )}
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
