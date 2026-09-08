import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useSession } from '../session/useSession'

export function ProtectedRoute() {
  const { session } = useSession()
  const location = useLocation()
  return session ? <Outlet /> : <Navigate to="/login" replace state={{ from: location }} />
}
