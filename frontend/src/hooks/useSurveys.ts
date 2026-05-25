import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import {
  createSurvey,
  getFeed,
  getSurvey,
  getSurveyResults,
  type FeedParams,
} from '@/services/surveys';
import type { CreateSurveyInput } from '@/types';

type CreateSurveyVariables = {
  data: CreateSurveyInput;
  access_token: string;
};

export function useFeed(params?: FeedParams) {
  return useQuery({
    queryKey: ['feed', params],
    queryFn: () => getFeed(params),
  });
}

export function useSurvey(id: string) {
  return useQuery({
    queryKey: ['survey', id],
    queryFn: () => getSurvey(id),
    enabled: id.length > 0,
  });
}

export function useSurveyResults(id: string) {
  return useQuery({
    queryKey: ['survey-results', id],
    queryFn: () => getSurveyResults(id),
    enabled: id.length > 0,
  });
}

export function useCreateSurvey() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ data, access_token }: CreateSurveyVariables) =>
      createSurvey(data, access_token),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['feed'] });
    },
  });
}
