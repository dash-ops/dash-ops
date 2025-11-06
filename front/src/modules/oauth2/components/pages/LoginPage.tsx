import { useLocation } from 'react-router';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
// import { Separator } from '@/components/ui/separator';
import { Github } from 'lucide-react';
// import { Chrome, LogIn } from 'lucide-react';
import { DarkModeToggle } from '@/components/theme/DarkModeToggle';

function LoginPage(): JSX.Element {
  const location = useLocation();
  const from = (location.state as string) || '/';
  const urlLoginGithub = `${
    (import.meta as { env?: Record<string, string> }).env?.VITE_API_URL ||
    'http://localhost:8080/api'
  }/oauth?redirect_url=${from}`;

  return (
    <div className="relative min-h-screen w-full flex items-center justify-center overflow-hidden bg-background">
      {/* Theme Toggle Button */}
      <div className="absolute top-4 right-4 z-20">
        <DarkModeToggle />
      </div>

      {/* Animated Grid Background */}
      <div className="absolute inset-0 overflow-hidden">
        {/* Floating Grid Lines */}
        <svg className="absolute inset-0 w-full h-full pointer-events-none">
          <defs>
            <pattern
              id="grid"
              width="40"
              height="40"
              patternUnits="userSpaceOnUse"
            >
              <path
                d="M 40 0 L 0 0 0 40"
                fill="none"
                stroke="currentColor"
                strokeWidth="0.5"
                className="text-border/80 dark:text-border/60"
              />
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#grid)" />
        </svg>

        {/* Animated Gradient Orbs */}
        <div className="absolute -top-1/4 -left-1/4 w-1/2 h-1/2 bg-primary/20 dark:bg-primary/10 rounded-full blur-3xl animate-pulse" />
        <div className="absolute -bottom-1/4 -right-1/4 w-1/2 h-1/2 bg-primary/20 dark:bg-primary/10 rounded-full blur-3xl animate-pulse" style={{ animationDelay: '1s' }} />
      </div>

      {/* Login Card */}
      <div className="relative z-10 w-full max-w-md px-4">
        <Card className="backdrop-blur-sm bg-card/95 border-border/50 shadow-2xl">
          <div className="p-8 space-y-6">
            {/* Logo and Title */}
            <div className="text-center space-y-2">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-primary/10 mb-4">
                <div className="animate-spin-slow">
                  <svg
                    width="32"
                    height="32"
                    viewBox="0 0 32 32"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <rect
                      x="4"
                      y="4"
                      width="10"
                      height="10"
                      rx="2"
                      className="fill-primary"
                    />
                    <rect
                      x="18"
                      y="4"
                      width="10"
                      height="10"
                      rx="2"
                      className="fill-primary/60"
                    />
                    <rect
                      x="4"
                      y="18"
                      width="10"
                      height="10"
                      rx="2"
                      className="fill-primary/60"
                    />
                    <rect
                      x="18"
                      y="18"
                      width="10"
                      height="10"
                      rx="2"
                      className="fill-primary/30"
                    />
                  </svg>
                </div>
              </div>

              <h1 className="text-foreground text-2xl font-bold">DashOPS</h1>
              <div className="flex items-center justify-center gap-2">
                <div className="h-px w-8 bg-border" />
                <p className="text-muted-foreground text-sm">Beta</p>
                <div className="h-px w-8 bg-border" />
              </div>
            </div>

            {/* Description */}
            <p className="text-center text-muted-foreground text-sm">
              Unified multi-cloud and Kubernetes management platform with integrated AI
            </p>

            {/* Login Buttons */}
            <div className="space-y-3">
              <div className="transition-transform hover:scale-[1.02] active:scale-[0.98]">
                <Button className="w-full gap-2 h-11" size="lg" asChild>
                  <a href={urlLoginGithub}>
                    <Github className="w-5 h-5" />
                    Continue with GitHub
                  </a>
                </Button>
              </div>

              {/* Other providers - commented until support is implemented */}
              {/* <div className="relative">
                <Separator className="my-4" />
                <span className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-card px-2 text-xs text-muted-foreground">
                  Other providers
                </span>
              </div>

              <div className="transition-transform hover:scale-[1.02] active:scale-[0.98]">
                <Button
                  className="w-full gap-2 h-11"
                  variant="outline"
                  disabled
                >
                  <Chrome className="w-5 h-5" />
                  Continue with Google
                </Button>
              </div>

              <div className="transition-transform hover:scale-[1.02] active:scale-[0.98]">
                <Button
                  className="w-full gap-2 h-11"
                  variant="outline"
                  disabled
                >
                  <LogIn className="w-5 h-5" />
                  Corporate SSO
                </Button>
              </div> */}
            </div>

            {/* Footer Info */}
            <div className="pt-4">
              <p className="text-xs text-center text-muted-foreground">
                By continuing, you agree to our Terms of Service
              </p>
            </div>
          </div>
        </Card>

        {/* Footer Links */}
        <div className="mt-4 flex items-center justify-center gap-4 text-xs text-muted-foreground">
          <a
            href="https://github.com/dash-ops/dash-ops"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-foreground transition-colors flex items-center gap-1"
            aria-label="View source code on GitHub"
          >
            <Github className="h-3 w-3" />
            <span>Source</span>
          </a>
          <span className="text-muted-foreground/50">•</span>
          <a
            href="https://dash-ops.github.io/"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-foreground transition-colors"
            aria-label="View documentation"
          >
            Docs
          </a>
        </div>
      </div>

      {/* Custom CSS for animations */}
      <style>{`
        @keyframes spin-slow {
          from {
            transform: rotate(0deg);
          }
          to {
            transform: rotate(360deg);
          }
        }
        .animate-spin-slow {
          animation: spin-slow 20s linear infinite;
        }
      `}</style>
    </div>
  );
}

export default LoginPage;
