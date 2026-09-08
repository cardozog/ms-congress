import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { useSession } from '../session/useSession'

export function AppLayout() {
  const { session, signOut } = useSession()
  const navigate = useNavigate()
  function logout() { signOut(); navigate('/login', { replace: true }) }
  return <div className="app-shell"><aside className="sidebar"><Brand /><div className="sidebar-rule" /><nav className="main-nav"><NavLink className={({ isActive }) => `nav-item${isActive ? ' is-active' : ''}`} to="/eventos">◈ <span>Eventos</span></NavLink><NavLink className={({ isActive }) => `nav-item${isActive ? ' is-active' : ''}`} to="/perfil">◎ <span>Perfil</span></NavLink></nav><div className="sidebar-footer"><span className="status-dot" /> API conectada</div></aside><main className="main-content"><header className="topbar"><div><p className="eyebrow">Área do organizador</p><h1>MS Congress</h1></div><div className="user-area"><div className="avatar">{session?.email.slice(0, 2).toUpperCase()}</div><div><strong>{session?.email}</strong><small>Pessoa jurídica</small></div><button className="icon-button" onClick={logout} title="Sair">↗</button></div></header><Outlet /></main></div>
}

function Brand() { return <a className="brand" href="/eventos"><span className="brand-mark">MS</span><span><strong>MS Congress</strong><small>painel do organizador</small></span></a> }
