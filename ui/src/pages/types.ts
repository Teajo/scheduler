export interface Task {
  date: string;
  publishers: Publisher[];
}

export interface Publisher {
  publisher: string;
  settings: Record<string, any>;
  retryStrategy: RetryStrategy;
}

export interface RetryStrategy {
  timeout: number;
  exponential: boolean;
  limit: number;
}
