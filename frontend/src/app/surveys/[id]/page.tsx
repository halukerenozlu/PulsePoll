'use client';

import { use, useState, useEffect } from 'react';
import Link from 'next/link';
import { useSurvey, useSurveyResults } from '@/hooks/useSurveys';
import { useVote, useChangeVote, useVerifyPin, useAcceptConsent, useReportSurvey } from '@/hooks/useVoting';
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

function ResultsSection({ id }: { id: string }) {
  const { data: results, isLoading, isError } = useSurveyResults(id);
  
  if (isLoading) return <div className="mt-8 p-4 bg-gray-50 rounded-xl animate-pulse text-sm text-gray-500">Sonuçlar yükleniyor...</div>;
  if (isError || !results) return <div className="mt-8 p-4 bg-red-50 text-red-600 rounded-xl text-sm">Sonuçlar yüklenemedi.</div>;

  return (
    <div className="mt-8 pt-6 border-t border-gray-200">
      <h3 className="text-lg font-medium text-gray-900 mb-4">Sonuçlar ({results.total_votes} toplam oy)</h3>
      <div className="space-y-4">
        {results.options.map(opt => (
          <div key={opt.id}>
            <div className="flex justify-between text-sm mb-1.5">
              <span className="text-gray-800 font-medium">{opt.text}</span>
              <span className="text-gray-600 font-medium">{opt.vote_count} oy ({opt.percentage}%)</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2.5 overflow-hidden">
              <div 
                className="bg-blue-500 h-full rounded-full transition-all duration-500 ease-out" 
                style={{ width: `${Math.max(0, Math.min(100, opt.percentage))}%` }}
              ></div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default function SurveyDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  
  const { data: survey, isLoading, isError } = useSurvey(id);
  
  // Consent state
  const [hasConsent, setHasConsent] = useState(true);
  const [dismissedConsent, setDismissedConsent] = useState(false);
  
  // PIN state
  const [pin, setPin] = useState('');
  const [pinVerified, setPinVerified] = useState(false);
  const [pinError, setPinError] = useState('');
  
  // Vote state
  const [isChangingVote, setIsChangingVote] = useState(false);
  const [voteError, setVoteError] = useState('');
  
  // Report state
  const [isReporting, setIsReporting] = useState(false);
  const [reportReason, setReportReason] = useState('');
  const [reportSuccess, setReportSuccess] = useState(false);

  const verifyPinMutation = useVerifyPin(id);
  const voteMutation = useVote(id);
  const changeVoteMutation = useChangeVote(id);
  const acceptConsentMutation = useAcceptConsent();
  const reportMutation = useReportSurvey(id);

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      const consent = localStorage.getItem('pulsepoll_consent');
      if (!consent) setHasConsent(false);
    }
  }, []);

  const handleVerifyPin = async (e: React.FormEvent) => {
    e.preventDefault();
    setPinError('');
    try {
      await verifyPinMutation.mutateAsync({ pin });
      setPinVerified(true);
    } catch (err: any) {
      if (err.message?.includes('429')) {
        setPinError('Çok fazla deneme. 15 dakika bekleyin.');
      } else {
        setPinError('Hatalı PIN girdiniz.');
      }
    }
  };

  const handleAcceptConsent = async () => {
    try {
      await acceptConsentMutation.mutateAsync();
      localStorage.setItem('pulsepoll_consent', 'true');
      setHasConsent(true);
      setVoteError(''); // Clear any previous consent errors
    } catch (err) {
      console.error('Consent failed', err);
    }
  };

  const handleVote = async (optionId: string) => {
    setVoteError('');
    const access_token = localStorage.getItem('access_token') || undefined;
    
    try {
      if (isChangingVote) {
        await changeVoteMutation.mutateAsync({ newOptionId: optionId, access_token });
        setIsChangingVote(false);
      } else {
        await voteMutation.mutateAsync({ optionId, access_token });
      }
    } catch (err: any) {
      const msg = err.message || '';
      if (msg.includes('PHASE_NOT_VOTING')) {
        setVoteError('Oylama süresi doldu');
      } else if (msg.includes('MAX_VOTES_REACHED')) {
        setVoteError('Maksimum oy hakkınızı kullandınız');
      } else if (msg.includes('VOTE_CHANGE_NOT_ALLOWED')) {
        setVoteError('Oy değiştirme hakkınızı kullandınız veya bu işlem izne tabi değil');
      } else if (msg.includes('429')) {
        setVoteError('Çok fazla istek. Lütfen bekleyin.');
      } else if (msg.includes('CONSENT_REQUIRED')) {
        setVoteError('Oy kullanabilmek için çerez izni vermeniz gerekmektedir.');
        setHasConsent(false);
        setDismissedConsent(false);
      } else {
        setVoteError('İşlem sırasında beklenmedik bir hata oluştu.');
      }
    }
  };

  const handleReport = async (e: React.FormEvent) => {
    e.preventDefault();
    const access_token = localStorage.getItem('access_token') || undefined;
    try {
      await reportMutation.mutateAsync({ reason: reportReason, access_token });
      setReportSuccess(true);
      setIsReporting(false);
    } catch (err) {
      alert('Rapor gönderilemedi. Lütfen daha sonra tekrar deneyin.');
    }
  };

  if (isLoading) {
    return (
      <div className="py-6 max-w-2xl mx-auto">
        <div className="mb-4">
          <Link href="/feed" className="text-sm font-medium text-gray-500 hover:text-black transition-colors">
            &larr; Akışa dön
          </Link>
        </div>
        <div className="p-6 md:p-8 bg-white border border-gray-200 rounded-2xl shadow-sm animate-pulse">
          <div className="h-7 bg-gray-200 rounded-md w-3/4 mb-4"></div>
          <div className="h-4 bg-gray-200 rounded-md w-full mb-2"></div>
          <div className="h-4 bg-gray-200 rounded-md w-5/6 mb-6"></div>
          <div className="space-y-3 mt-8">
            {[1, 2, 3].map(i => (
              <div key={i} className="h-14 bg-gray-100 rounded-xl w-full"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (isError || !survey) {
    return (
      <div className="py-12 max-w-xl mx-auto text-center">
        <div className="p-10 bg-white border border-red-100 rounded-2xl shadow-sm">
          <h2 className="text-xl font-bold text-gray-900 mb-2">Anket bulunamadı</h2>
          <p className="text-sm text-gray-500 mb-8">Aradığınız anket silinmiş veya mevcut olmayabilir.</p>
          <Link href="/feed" className="px-5 py-2.5 bg-black text-white font-medium rounded-lg hover:bg-gray-800 transition-colors">
            Akışa Dön
          </Link>
        </div>
      </div>
    );
  }

  const showPinForm = survey.requires_pin && !pinVerified;
  const showConsentBanner = survey.can_vote && !showPinForm && !hasConsent && !dismissedConsent;
  const showVotingOptions = survey.can_vote && survey.phase === 'VOTING' && !showPinForm;

  return (
    <div className="py-6 max-w-2xl mx-auto">
      <div className="mb-6">
        <Link href="/feed" className="text-sm font-medium text-gray-500 hover:text-black transition-colors">
          &larr; Akışa dön
        </Link>
      </div>

      <div className="p-6 md:p-8 bg-white border border-gray-200 rounded-2xl shadow-sm">
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-900 leading-tight mb-3">
            {survey.title}
          </h1>
          {survey.description && (
            <p className="text-gray-600 text-base mb-5 whitespace-pre-wrap leading-relaxed">
              {survey.description}
            </p>
          )}
          <div className="flex flex-wrap items-center gap-3 mt-2">
            <PhaseBadge phase={survey.phase} />
            {survey.phase === 'VOTING' && survey.vote_ends_at && (
              <span className="text-sm text-gray-500 font-medium">
                {getTimeRemainingText(survey.vote_ends_at)}
              </span>
            )}
          </div>
        </div>

        {survey.phase === 'EXPIRED' && (
          <div className="mb-6 p-4 bg-gray-50 border border-gray-200 rounded-xl text-center">
            <p className="text-gray-700 font-medium">Bu anketin süresi doldu.</p>
          </div>
        )}

        {showPinForm && (
          <div className="mt-8 p-6 border border-gray-200 rounded-xl bg-gray-50">
            <h3 className="text-lg font-bold text-gray-900 mb-2">PIN Korumalı Anket</h3>
            <p className="text-sm text-gray-600 mb-5">
              Oy kullanabilmek veya sonuçları görebilmek için lütfen anket sahibinden aldığınız PIN kodunu girin.
            </p>
            {pinError && (
              <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-md border border-red-100">
                {pinError}
              </div>
            )}
            <form onSubmit={handleVerifyPin} className="flex gap-3">
              <input 
                type="text" 
                value={pin}
                onChange={e => setPin(e.target.value)}
                placeholder="PIN Kodu"
                required
                className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black text-gray-900"
              />
              <button 
                type="submit"
                disabled={verifyPinMutation.isPending}
                className="px-5 py-2 bg-black text-white font-medium rounded-lg hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {verifyPinMutation.isPending ? 'Doğrulanıyor...' : 'Doğrula'}
              </button>
            </form>
          </div>
        )}

        {showConsentBanner && (
          <div className="mt-6 p-5 border border-blue-200 bg-blue-50/50 rounded-xl">
            <h3 className="text-base font-bold text-blue-900 mb-2">Oy kullanmak için gerekli çerez</h3>
            <p className="text-sm text-blue-800 mb-5 leading-relaxed">
              Oyları adil tutmak, spam’i engellemek ve "1 kez oy değiştirme" hakkını uygulamak için gerekli bir çerez kullanıyoruz.
              Kişisel bilgi içermez ve 48 saat içinde otomatik olarak geçersiz olur. Kabul etmeden de anketleri ve sonuçları görüntüleyebilirsin.
            </p>
            <div className="flex flex-wrap items-center gap-3">
              <button 
                onClick={handleAcceptConsent}
                disabled={acceptConsentMutation.isPending}
                className="px-5 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors shadow-sm"
              >
                {acceptConsentMutation.isPending ? 'Bekleyin...' : 'Kabul et & Oy ver'}
              </button>
              <button 
                onClick={() => setDismissedConsent(true)}
                className="px-5 py-2 bg-white text-blue-700 border border-blue-200 text-sm font-medium rounded-lg hover:bg-blue-50 hover:border-blue-300 transition-colors shadow-sm"
              >
                Şimdi değil (sadece görüntüle)
              </button>
            </div>
          </div>
        )}

        {showVotingOptions && (
          <div className="mt-8">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Seçenekler</h3>
            {voteError && (
              <div className="mb-5 p-3 bg-red-50 text-red-700 text-sm rounded-md border border-red-100">
                {voteError}
              </div>
            )}
            <div className="space-y-3">
              {survey.options.map(opt => (
                <div 
                  key={opt.id} 
                  className="flex items-center justify-between p-4 border border-gray-200 rounded-xl hover:border-gray-300 hover:bg-gray-50 transition-all group"
                >
                  <span className="text-gray-900 font-medium">{opt.text}</span>
                  <button 
                    onClick={() => handleVote(opt.id)}
                    disabled={voteMutation.isPending || changeVoteMutation.isPending || (survey.requires_pin && !pinVerified)}
                    className={`px-5 py-2 text-sm font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed ${
                      isChangingVote 
                        ? 'bg-blue-600 hover:bg-blue-700 text-white' 
                        : 'bg-black hover:bg-gray-800 text-white'
                    }`}
                  >
                    {isChangingVote ? 'Buna Değiştir' : 'Oy Ver'}
                  </button>
                </div>
              ))}
            </div>
            
            {survey.allow_vote_change_once && !isChangingVote && (
              <div className="mt-5 text-right">
                <button 
                  onClick={() => setIsChangingVote(true)} 
                  className="text-sm font-medium text-gray-500 hover:text-black underline underline-offset-4 transition-colors"
                >
                  Oyumu değiştirmek istiyorum
                </button>
              </div>
            )}
            {isChangingVote && (
              <div className="mt-5 text-right">
                <button 
                  onClick={() => setIsChangingVote(false)} 
                  className="text-sm font-medium text-gray-500 hover:text-black underline underline-offset-4 transition-colors"
                >
                  İptal et
                </button>
              </div>
            )}
          </div>
        )}

        {survey.results_visible && !showPinForm && (
          <ResultsSection id={survey.id} />
        )}
      </div>

      <div className="mt-8 mb-12 text-center">
        {!isReporting ? (
          <button 
            onClick={() => setIsReporting(true)}
            className="text-sm font-medium text-gray-400 hover:text-gray-600 transition-colors"
          >
            Bu anketi raporla
          </button>
        ) : reportSuccess ? (
          <div className="inline-block px-4 py-2 bg-green-50 text-green-700 text-sm font-medium rounded-lg border border-green-100">
            Raporunuz başarıyla gönderildi.
          </div>
        ) : (
          <form 
            onSubmit={handleReport} 
            className="max-w-sm mx-auto p-5 bg-white border border-gray-200 rounded-xl shadow-sm text-left"
          >
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Rapor Sebebi
            </label>
            <input 
              type="text" 
              value={reportReason}
              onChange={e => setReportReason(e.target.value)}
              placeholder="Örn: Uygunsuz içerik, spam..."
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black mb-4 text-sm text-gray-900"
            />
            <div className="flex justify-end gap-3">
              <button 
                type="button" 
                onClick={() => setIsReporting(false)}
                className="px-4 py-2 text-sm font-medium text-gray-600 hover:text-black transition-colors"
              >
                İptal
              </button>
              <button 
                type="submit"
                disabled={reportMutation.isPending}
                className="px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-lg hover:bg-red-700 disabled:opacity-50 transition-colors"
              >
                {reportMutation.isPending ? 'Gönderiliyor...' : 'Gönder'}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
