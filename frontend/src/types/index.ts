export type AuthUser = {
  id: string;
  email: string;
  display_name: string;
};

export type AuthResponse = {
  access_token: string;
  user: AuthUser;
};

export type LogoutResponse = {
  ok: boolean;
};

export type SurveyVisibility = 'public' | 'unlisted' | 'private_pin';

export type SurveyResultsMode = 'open_live' | 'closed_hidden_until_end';

export type SurveyPhase = 'VOTING' | 'RESULTS' | 'EXPIRED';

export type SurveyFeedItem = {
  id: string;
  title: string;
  description?: string;
  visibility: SurveyVisibility;
  results_mode: SurveyResultsMode;
  created_at: string;
  vote_ends_at: string;
  results_ends_at: string;
  phase: SurveyPhase;
  can_vote: boolean;
  results_visible: boolean;
  requires_pin: boolean;
};

export type SurveyOption = {
  id: string;
  text: string;
  position: number;
};

export type Survey = SurveyFeedItem & {
  creator_id: string;
  max_votes_per_user: number;
  allow_vote_change_once: boolean;
  retention_ends_at: string;
  options: SurveyOption[];
};

export type CreateSurveyInput = {
  title: string;
  description?: string;
  options: string[];
  visibility: SurveyVisibility;
  access_pin?: string;
  results_mode: SurveyResultsMode;
  max_votes_per_user?: number;
  allow_vote_change_once?: boolean;
  vote_ends_at?: string;
  results_ends_at?: string;
  retention_ends_at?: string;
};

export type SurveyResultOption = {
  id: string;
  text: string;
  vote_count: number;
  percentage: number;
};

export type SurveyResults = {
  survey_id: string;
  total_votes: number;
  options: SurveyResultOption[];
};
