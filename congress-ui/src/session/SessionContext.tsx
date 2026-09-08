import { useMemo, useState, type ReactNode } from 'react'
import { clearSession, saveSession } from '../services/api'
import type { Session } from '../types'
import { SessionContext } from './context'

export function SessionProvider({ initialSession, children }: { initialSession: Session | null; children: ReactNode }) {
  const [session, setSession] = useState<Session | null>(initialSession)
  const value = useMemo(() => ({
    session,
    signIn: (next: Session) => { saveSession(next); setSession(next) },
    signOut: () => { clearSession(); setSession(null) },
  }), [session])
  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>
}

