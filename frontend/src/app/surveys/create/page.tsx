'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useCreateSurvey } from '@/hooks/useSurveys';
import type { SurveyVisibility, SurveyResultsMode } from '@/types';

export default function CreateSurveyPage() {
  const router = useRouter();
  const createMutation = useCreateSurvey();

  // Authentication check
  const [token, setToken] = useState<string | null>(null);
  const [isCheckingAuth, setIsCheckingAuth] = useState(true);

  // Form states
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [options, setOptions] = useState<string[]>(['', '']);
  const [visibility, setVisibility] = useState<SurveyVisibility>('public');
  const [accessPin, setAccessPin] = useState('');
  const [resultsMode, setResultsMode] = useState<SurveyResultsMode>('open_live');
  const [maxVotes, setMaxVotes] = useState<number>(1);
  const [allowChange, setAllowChange] = useState<boolean>(false);

  const [errorMessage, setErrorMessage] = useState('');

  useEffect(() => {
    const accessToken = localStorage.getItem('access_token');
    if (!accessToken) {
      router.push('/login');
    } else {
      setToken(accessToken);
      setIsCheckingAuth(false);
    }
  }, [router]);

  const handleOptionChange = (index: number, value: string) => {
    const newOptions = [...options];
    newOptions[index] = value;
    setOptions(newOptions);
  };

  const handleAddOption = () => {
    if (options.length < 10) {
      setOptions([...options, '']);
    }
  };

  const handleRemoveOption = (index: number) => {
    if (options.length > 2) {
      const newOptions = [...options];
      newOptions.splice(index, 1);
      setOptions(newOptions);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMessage('');

    if (!token) {
      router.push('/login');
      return;
    }

    const filteredOptions = options.map(o => o.trim()).filter(o => o.length > 0);
    if (filteredOptions.length < 2) {
      setErrorMessage('Lütfen en az 2 geçerli seçenek girin.');
      return;
    }

    try {
      const payload = {
        title,
        description: description.trim() ? description : undefined,
        options: filteredOptions,
        visibility,
        access_pin: visibility === 'private_pin' ? accessPin : undefined,
        results_mode: resultsMode,
        max_votes_per_user: maxVotes,
        allow_vote_change_once: maxVotes === 1 ? allowChange : undefined,
      };

      const survey = await createMutation.mutateAsync({ data: payload, access_token: token });
      router.push(`/surveys/${survey.id}`);
    } catch (error: any) {
      const msg = error.message || '';
      if (msg.includes('401')) {
        localStorage.removeItem('access_token');
        router.push('/login');
      } else if (msg.includes('403') || msg.toLowerCase().includes('moderation') || msg.toLowerCase().includes('blocked')) {
        setErrorMessage('Anket içeriği uygunsuz kelimeler içeriyor. Lütfen başlık, açıklama ve seçenekleri kontrol edin.');
      } else if (msg.includes('400')) {
        setErrorMessage('Girdiğiniz bilgiler geçersiz. Lütfen alanları kontrol edip tekrar deneyin.');
      } else {
        setErrorMessage('Anket oluşturulurken beklenmedik bir hata meydana geldi. Lütfen daha sonra tekrar deneyin.');
      }
    }
  };

  if (isCheckingAuth) {
    return (
      <div className="py-12 text-center text-gray-500">
        Yetkilendirme kontrol ediliyor...
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto py-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Yeni Anket Oluştur</h1>
        <p className="text-sm text-gray-500 mt-1">Düşünceleri hızlıca öğrenin</p>
      </div>

      {errorMessage && (
        <div className="mb-6 p-4 bg-red-50 text-red-700 text-sm rounded-lg border border-red-100">
          {errorMessage}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-8 bg-white p-6 md:p-8 rounded-2xl border border-gray-200 shadow-sm">
        {/* Title */}
        <div>
          <label htmlFor="title" className="block text-sm font-medium text-gray-900 mb-1.5">
            Anket Başlığı <span className="text-red-500">*</span>
          </label>
          <input
            id="title"
            type="text"
            required
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Ne sormak istersiniz?"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black text-gray-900"
          />
        </div>

        {/* Description */}
        <div>
          <label htmlFor="description" className="block text-sm font-medium text-gray-900 mb-1.5">
            Açıklama (İsteğe bağlı)
          </label>
          <textarea
            id="description"
            rows={3}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Anketinizle ilgili ek detaylar..."
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black text-gray-900 resize-y"
          />
        </div>

        {/* Options */}
        <div className="space-y-3">
          <label className="block text-sm font-medium text-gray-900">
            Seçenekler <span className="text-red-500">*</span>
          </label>
          
          {options.map((option, index) => (
            <div key={index} className="flex gap-2 items-center">
              <input
                type="text"
                required
                value={option}
                onChange={(e) => handleOptionChange(index, e.target.value)}
                placeholder={`Seçenek ${index + 1}`}
                className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black text-gray-900"
              />
              {options.length > 2 && (
                <button
                  type="button"
                  onClick={() => handleRemoveOption(index)}
                  className="p-2 text-gray-400 hover:text-red-600 focus:outline-none transition-colors"
                  aria-label="Seçeneği Kaldır"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                  </svg>
                </button>
              )}
            </div>
          ))}

          {options.length < 10 && (
            <button
              type="button"
              onClick={handleAddOption}
              className="mt-2 text-sm font-medium text-blue-600 hover:text-blue-800 transition-colors"
            >
              + Yeni Seçenek Ekle
            </button>
          )}
        </div>

        <hr className="border-gray-200" />

        {/* Visibility */}
        <div>
          <label className="block text-sm font-medium text-gray-900 mb-3">
            Görünürlük <span className="text-red-500">*</span>
          </label>
          <div className="space-y-2">
            <label className="flex items-center gap-3">
              <input
                type="radio"
                name="visibility"
                value="public"
                checked={visibility === 'public'}
                onChange={() => setVisibility('public')}
                className="w-4 h-4 text-black focus:ring-black border-gray-300"
              />
              <span className="text-sm text-gray-800">Herkese Açık (Akışta görünür)</span>
            </label>
            <label className="flex items-center gap-3">
              <input
                type="radio"
                name="visibility"
                value="unlisted"
                checked={visibility === 'unlisted'}
                onChange={() => setVisibility('unlisted')}
                className="w-4 h-4 text-black focus:ring-black border-gray-300"
              />
              <span className="text-sm text-gray-800">Gizli Link (Sadece link ile erişilir)</span>
            </label>
            <label className="flex items-center gap-3">
              <input
                type="radio"
                name="visibility"
                value="private_pin"
                checked={visibility === 'private_pin'}
                onChange={() => setVisibility('private_pin')}
                className="w-4 h-4 text-black focus:ring-black border-gray-300"
              />
              <span className="text-sm text-gray-800">PIN Korumalı (Link + PIN kodu gerekir)</span>
            </label>
          </div>
        </div>

        {visibility === 'private_pin' && (
          <div className="bg-gray-50 p-4 rounded-xl border border-gray-200">
            <label htmlFor="access_pin" className="block text-sm font-medium text-gray-900 mb-1.5">
              Erişim PIN Kodu <span className="text-red-500">*</span>
            </label>
            <input
              id="access_pin"
              type="text"
              required
              value={accessPin}
              onChange={(e) => setAccessPin(e.target.value)}
              placeholder="Örn: 123456"
              className="w-full md:w-1/2 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black text-gray-900"
            />
          </div>
        )}

        {/* Results Mode */}
        <div>
          <label className="block text-sm font-medium text-gray-900 mb-3">
            Sonuçların Görünürlüğü <span className="text-red-500">*</span>
          </label>
          <div className="space-y-2">
            <label className="flex items-center gap-3">
              <input
                type="radio"
                name="resultsMode"
                value="open_live"
                checked={resultsMode === 'open_live'}
                onChange={() => setResultsMode('open_live')}
                className="w-4 h-4 text-black focus:ring-black border-gray-300"
              />
              <span className="text-sm text-gray-800">Oylar anlık görünsün</span>
            </label>
            <label className="flex items-center gap-3">
              <input
                type="radio"
                name="resultsMode"
                value="closed_hidden_until_end"
                checked={resultsMode === 'closed_hidden_until_end'}
                onChange={() => setResultsMode('closed_hidden_until_end')}
                className="w-4 h-4 text-black focus:ring-black border-gray-300"
              />
              <span className="text-sm text-gray-800">Sadece oylama bitince görünsün</span>
            </label>
          </div>
        </div>

        {/* Voting Limits */}
        <div>
          <label htmlFor="maxVotes" className="block text-sm font-medium text-gray-900 mb-1.5">
            Kullanıcı Başına Maksimum Oy
          </label>
          <input
            id="maxVotes"
            type="number"
            min={1}
            value={maxVotes}
            onChange={(e) => setMaxVotes(parseInt(e.target.value) || 1)}
            className="w-full md:w-1/3 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-black text-gray-900"
          />
        </div>

        {maxVotes === 1 && (
          <div>
            <label className="flex items-center gap-3">
              <input
                type="checkbox"
                checked={allowChange}
                onChange={(e) => setAllowChange(e.target.checked)}
                className="w-4 h-4 text-black focus:ring-black border-gray-300 rounded"
              />
              <span className="text-sm text-gray-900 font-medium">1 kez oy değiştirme hakkı ver</span>
            </label>
          </div>
        )}

        <div className="pt-4 border-t border-gray-200">
          <button
            type="submit"
            disabled={createMutation.isPending}
            className="w-full md:w-auto px-8 py-3 bg-black text-white font-medium rounded-lg hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-black disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {createMutation.isPending ? 'Oluşturuluyor...' : 'Anketi Oluştur'}
          </button>
        </div>
      </form>
    </div>
  );
}
