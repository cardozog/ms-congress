import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useSession } from '../session/useSession'

export function ProtectedRoute({ role }: { role?: string | string[] }) {
  const { session } = useSession()
  const location = useLocation()
  const requiredRoles = role ? (Array.isArray(role) ? role : [role]) : []
  if (!session) return <Navigate to={requiredRoles.includes('PARTICIPANTE') ? '/congressista/login' : '/login'} replace state={{ from: location }} />
  if (requiredRoles.length > 0 && !requiredRoles.some(requiredRole => session.papeis?.includes(requiredRole))) return <Navigate to={session.papeis?.includes('PARTICIPANTE') ? '/congressista/eventos' : '/eventos'} replace />
  return <Outlet />
}
