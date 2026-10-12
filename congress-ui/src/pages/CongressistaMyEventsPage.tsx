import { useEffect, useMemo, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getMyEvents } from '../services/api'
import { useSession } from '../session/useSession'
import type { MyEvent } from '../types'
import { errorMessage } from '../utils'

type Filter = 'UPCOMING' | 'PREVIOUS' | 'ALL'
const date = (value: string) => new Date(value).toLocaleDateString('pt-BR', { day: '2-digit', month: 'long', year: 'numeric' })

export function CongressistaMyEventsPage() {
  const { session } = useSession()
  const navigate = useNavigate()
  const [events, setEvents] = useState<MyEvent[]>([])
  const [filter, setFilter] = useState<Filter>('UPCOMING')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  useEffect(() => {
    if (!session) return
    void getMyEvents(session).then(setEvents).catch(error => setError(errorMessage(error, 'Não foi possível carregar seus eventos.'))).finally(() => setLoading(false))
  }, [session])
  const visible = useMemo(() => events.filter(event => filter === 'ALL' || (filter === 'UPCOMING' ? event.classificacao !== 'ENDED' : event.classificacao === 'ENDED')), [events, filter])
  return <div className="app-shell congressista-shell"><main className="main-content"><header className="public-nav"><Link className="public-brand" to="/congressista/eventos"><span className="brand-mark">MS</span><strong>MS Congress</strong></Link><nav><Link to="/congressista/eventos">Explorar eventos</Link><Link className="is-current" to="/meus-eventos">Meus Eventos</Link></nav><div className="public-actions"><span className="public-user">{session?.email}</span></div></header><section className="my-events-view"><nav className="event-filters" aria-label="Filtrar eventos"><button className={filter === 'UPCOMING' ? 'active' : ''} onClick={() => setFilter('UPCOMING')}>Próximos</button><button className={filter === 'PREVIOUS' ? 'active' : ''} onClick={() => setFilter('PREVIOUS')}>Anteriores</button><button className={filter === 'ALL' ? 'active' : ''} onClick={() => setFilter('ALL')}>Todos</button></nav>{loading ? <div className="empty-state"><p>Carregando seus eventos...</p></div> : error ? <div className="empty-state"><h3>Não foi possível carregar</h3><p>{error}</p></div> : !visible.length ? <div className="empty-state"><h3>{events.length ? 'Nenhum evento neste filtro' : 'Nenhum ingresso adquirido'}</h3><p>{events.length ? 'Tente outro filtro para encontrar seus eventos.' : 'Após um pagamento aprovado, seus eventos aparecerão aqui.'}</p></div> : <div className="my-event-list">{visible.map(event => <article className="my-event-card" key={event.id}><div className="my-event-cover">{event.eventoLogo ? <img src={event.eventoLogo} alt="" /> : <span>MS Congress</span>}</div><div className="my-event-content"><span className="event-status">{event.classificacao === 'ONGOING' ? 'Em andamento' : event.classificacao === 'ENDED' ? 'Encerrado' : 'Próximo'}</span><h2>{event.nome}</h2><p>{date(event.dataInicio)} — {date(event.dataFim)}</p><strong>🎟 {event.quantidade} ingresso{event.quantidade === 1 ? '' : 's'}</strong><span className="purchase-confirmed">Compra confirmada</span><button className="primary-button" type="button" onClick={() => navigate(`/meus-eventos/${event.id}`)}>Ver ingressos</button></div></article>)}</div>}</section></main></div>
}
