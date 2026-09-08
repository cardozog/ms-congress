import { createContext } from 'react'
import type { Session } from '../types'

export interface SessionContextValue {
  session: Session | null
  signIn: (session: Session) => void
  signOut: () => void
}

export const SessionContext = createContext<SessionContextValue | null>(null)