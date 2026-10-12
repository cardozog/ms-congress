import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getMyEvent } from '../services/api'
import { useSession } from '../session/useSession'
import type { MyEvent } from '../types'
import { errorMessage } from '../utils'

export function CongressistaMyEventDetailPage() {
  const { session } = useSession()
  const { id } = useParams()
  const [event, setEvent] = useState<MyEvent | null>(null)
  const [error, setError] = useState('')
  useEffect(() => { if (session && id) void getMyEvent(session, Number(id)).then(setEvent).catch(error => setError(errorMessage(error, 'Não foi possível carregar os ingressos.'))) }, [session, id])
  return <div className="app-shell congressista-shell"><main className="main-content"><section className="my-events-view">  <Link className="secondary-button" to="/meus-eventos">← Meus Eventos</Link>{error ? <div className="empty-state"><p>{error}</p></div> : !event ? <div className="empty-state"><p>Carregando ingressos...</p></div> : <><span className="kicker">Credenciais digitais</span><h1>{event.nome}</h1><p>{event.quantidade} ingressos adquiridos</p><div className="credential-list">{event.ingressos.map(ticket => <div className="credential-card" key={ticket.id}><strong>{ticket.ingresso}</strong><span>Compra confirmada</span><code>{ticket.codigo}</code></div>)}</div></>}</section></main></div>
}
