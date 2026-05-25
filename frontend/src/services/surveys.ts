import { api } from './api';
import type { CreateSurveyInput, Survey, SurveyFeedItem, SurveyResults } from '@/types';

export type FeedParams = {
  sort?: 'new';
  search?: string;
};

export type FeedResponse = {
  items: SurveyFeedItem[];
};

export function getFeed(params: FeedParams = {}): Promise<FeedResponse> {
  const searchParams = new URLSearchParams();

  if (params.sort) {
    searchParams.set('sort', params.sort);
  }
  if (params.search) {
    searchParams.set('search', params.search);
  }

  const query = searchParams.toString();
  return api<FeedResponse>(`/api/v1/feed${query ? `?${query}` : ''}`);
}

export function getSurvey(id: string): Promise<Survey> {
  return api<Survey>(`/api/v1/surveys/${id}`);
}

export function createSurvey(data: CreateSurveyInput, access_token: string): Promise<Survey> {
  return api<Survey>('/api/v1/surveys', {
    method: 'POST',
    headers: { Authorization: `Bearer ${access_token}` },
    body: JSON.stringify(data),
  });
}

export function getSurveyResults(id: string): Promise<SurveyResults> {
  return api<SurveyResults>(`/api/v1/surveys/${id}/results`);
}
