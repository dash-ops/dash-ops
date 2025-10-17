import { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { toast } from 'sonner';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Cloud } from 'lucide-react';
import { AWSTypes } from '@/types';
import { Page } from '@/types';
import ContentWithMenu from '../../../components/ContentWithMenu';
import { getAccountsCached } from '../utils/accountsCache';

interface AWSWithAccountSelectorProps {
  pages: Page[];
}

export default function AWSWithAccountSelector({
  pages,
}: AWSWithAccountSelectorProps): JSX.Element {
  const location = useLocation();
  const navigate = useNavigate();
  const [accounts, setAccounts] = useState<AWSTypes.Account[]>([]);
  const [selectedAccountKey, setSelectedAccountKey] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [isLoadingAccounts, setIsLoadingAccounts] = useState(false);

  useEffect(() => {
    if (isLoadingAccounts || accounts.length > 0) return;

    async function loadAccounts() {
      setIsLoadingAccounts(true);
      try {
        const data = await getAccountsCached();
        setAccounts(data);
        setLoading(false);
      } catch (error) {
        console.error('Failed to load accounts:', error);
        toast.error('Failed to load AWS accounts');
        setLoading(false);
      } finally {
        setIsLoadingAccounts(false);
      }
    }

    loadAccounts();
  }, [accounts.length, isLoadingAccounts]);

  useEffect(() => {
    if (accounts.length === 0) return;

    const pathParts = location.pathname.split('/');
    const keyIndex = pathParts.findIndex((part) => part === 'aws') + 1;
    const currentKey = pathParts[keyIndex];

    if (currentKey && accounts.some((account) => account.key === currentKey)) {
      setSelectedAccountKey(currentKey);
    } else {
      setSelectedAccountKey('');
    }
  }, [location.pathname, accounts]);

  const handleAccountChange = (newAccountKey: string) => {
    setSelectedAccountKey(newAccountKey);
    navigate(`/aws/${newAccountKey}`);
  };

  const pagesWithContext = pages;

  if (!loading && accounts.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 space-y-4">
        <Cloud className="h-12 w-12 text-muted-foreground" />
        <div className="text-center space-y-2">
          <h3 className="font-semibold text-lg">No AWS Accounts</h3>
          <p className="text-muted-foreground max-w-md">
            No AWS accounts are configured. Please check your configuration
            file.
          </p>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
        <span className="ml-3">Loading accounts...</span>
      </div>
    );
  }

  const selectedAccount = accounts.find(
    (account) => account.key === selectedAccountKey
  );

  return (
    <div className="space-y-6">
      <div className="border-b pb-4">
        <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
          <div className="space-y-1">
            <h1 className="text-2xl font-semibold tracking-tight">AWS</h1>
            <p className="text-muted-foreground">
              Manage your AWS accounts and resources
            </p>
          </div>

          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Cloud className="h-4 w-4" />
              Account:
            </div>
            <Select
              value={selectedAccountKey}
              onValueChange={handleAccountChange}
            >
              <SelectTrigger className="w-[280px]">
                <SelectValue placeholder="Select an account" />
              </SelectTrigger>
              <SelectContent>
                {accounts.map((account) => (
                  <SelectItem key={account.key} value={account.key}>
                    <div className="flex items-center gap-2">
                      <span>{account.name}</span>
                      <Badge variant="secondary" className="text-xs">
                        {account.key}
                      </Badge>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
        {selectedAccount && (
          <div className="flex items-center gap-2 text-sm text-muted-foreground mt-2">
            <span>Selected:</span>
            <Badge variant="outline">
              {selectedAccount.name} ({selectedAccount.key})
            </Badge>
          </div>
        )}
      </div>

      {!selectedAccountKey ? (
        <div className="text-center py-8">
          <span className="text-muted-foreground">No AWS account provided</span>
        </div>
      ) : (
        <ContentWithMenu
          pages={pagesWithContext}
          paramName="key"
          contextValue={selectedAccountKey}
        />
      )}
    </div>
  );
}
