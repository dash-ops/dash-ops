export type DataSource = 'Logs' | 'Traces' | 'Metrics';

export interface QueryTab {
  id: string;
  title: string;
  query: string;
  dataSource: DataSource | null;
  results: any[] | null;
  isActive: boolean;
  timestamp: number;
}

