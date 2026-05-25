import Providers from './providers';
import { Navbar } from '@/components/Navbar';
import './globals.css';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <Providers>
          <div className="min-h-screen flex flex-col">
            <Navbar />
            <main className="flex-1 w-full max-w-4xl mx-auto p-4 md:p-6">
              {children}
            </main>
          </div>
        </Providers>
      </body>
    </html>
  );
}
