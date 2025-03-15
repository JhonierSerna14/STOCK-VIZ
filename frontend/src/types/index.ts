export interface Stock {
  ticker: string;
  company: string;
  brokerage: string;
  action: string;
  rating_from: string;
  rating_to: string;
  target_from: string;
  target_to: string;
  time: string;
}

export interface Pagination {
  current_page: number;
  per_page: number;
  total_items: number;
  total_pages: number;
}

export interface StockResponse {
  items: Stock[];
  pagination: Pagination;
}

export interface Recommendation {
  ticker: string;
  company: string;
  score: number;
  latest_rating: string;
  latest_target: string;
  last_updated: string;
  analysis_rationale: string;
}

export interface RecommendationFilter {
  limit?: number;
  date_from?: string;
  date_to?: string;
  rating?: string;
  ticker?: string;
}
