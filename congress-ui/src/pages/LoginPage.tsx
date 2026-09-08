import { Navigate, useNavigate } from 'react-router-dom'
import { LoginScreen } from '../components/LoginScreen'
import { useSession } from '../session/useSession'
import type { Session } from '../types'

export function LoginPage() {
  const { session, signIn } = useSession()
  const navigate = useNavigate()
  if (session) return <Navigate to="/eventos" replace />
  async function handleLogin(next: Session) { signIn(next); navigate('/eventos', { replace: true }) }
  return <><header className="topbar public-topbar"><Brand /><div /></header><LoginScreen onLogin={handleLogin} /></>
}
function Brand() { return <a className="brand" href="/login"><span className="brand-mark">MS</span><span><strong>MS Congress</strong><small>painel do organizador</small></span></a> }
