import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getAvailableEvents, getCart } from '../services/api'
import { useSession } from '../session/useSession'
import type { AvailableEventGroup } from '../types'
import { errorMessage } from '../utils'

export function CongressistaEventsPage() {
  const { session, signOut } = useSession()
  const navigate = useNavigate()
  const [groups, setGroups] = useState<AvailableEventGroup[]>([])
  const [loading, setLoading] = useState(true)
  const [feedback, setFeedback] = useState('')
  const [cartCount, setCartCount] = useState(0)
  useEffect(() => {
    if (!session) return
    const currentSession = session
    async function load() {
      setLoading(true); setFeedback('')
      try { const [available, cart] = await Promise.all([getAvailableEvents(currentSession), getCart(currentSession)]); setGroups(available); setCartCount(cart.totalItens) }
      catch (error) { setFeedback(errorMessage(error, 'Não foi possível carregar os eventos disponíveis.')) }
      finally { setLoading(false) }
    }
    void load()
  }, [session])
  function logout() { signOut(); navigate('/congressista/login', { replace: true }) }
  return <div className="app-shell congressista-shell"><main className="main-content"><header className="topbar"><div><p className="eyebrow">Área do congressista</p><h1>Descobrir eventos</h1></div><div className="user-area">  <Link className="secondary-button" to="/meus-eventos">Meus Eventos</Link><Link className="cart-link" to="/congressista/carrinho">Carrinho <b>{cartCount}</b></Link><div><strong>{session?.email}</strong><small>Congressista</small></div><button className="icon-button" onClick={logout} title="Sair">↗</button></div></header><section className="discovery-view"><div className="discovery-heading"><div><span className="kicker">Programação aberta</span><h2>Eventos disponíveis</h2><p>Encontre experiências para participar hoje.</p></div><span className="discovery-search" aria-hidden="true">⌕</span></div>{loading ? <div className="empty-state"><p>Carregando eventos...</p></div> : feedback ? <div className="empty-state"><h3>Não foi possível carregar</h3><p>{feedback}</p></div> : groups.length === 0 ? <div className="empty-state"><h3>Nenhum evento disponível</h3><p>Novos eventos publicados aparecerão aqui.</p></div> : <div className="organizer-list">{groups.map(group => <section className="organizer-section" key={group.organizador.id}><h3>{group.organizador.nomeFantasia || group.organizador.razaoSocial}</h3><div className="event-rail">{group.eventos.map(event => <button className="discovery-card" type="button" key={event.id} onClick={() => navigate(`/congressista/eventos/${event.id}`)}><span className="discovery-image">{event.eventoLogo ? <img src={event.eventoLogo} alt="" /> : <span>MS Congress</span>}</span><span className="discovery-overlay"><strong>{event.nome}</strong></span><span className="discovery-card-title">{event.nome}</span></button>)}</div></section>)}</div>}</section></main></div>
}
