import type { EventSummary } from '../types'
import { formatDate } from '../utils'
import { CircleStop, Send } from 'lucide-react'

export function EventCard({ event, onSelect, onPublish, onClose, publishing, closing }: { event: EventSummary; onSelect: (id: number) => void; onPublish: (id: number) => void; onClose: (id: number) => void; publishing: boolean; closing: boolean }) {
  const status = event.encerrado ? 'Encerrado' : event.publicado ? 'Publicado' : 'Rascunho'
  return <article className="event-card"><button className="event-card-content" type="button" onClick={() => onSelect(event.id)}><span className="event-date">{formatDate(event.dataInicio)}</span><span className={`event-status ${event.publicado ? 'is-published' : ''} ${event.encerrado ? 'is-closed' : ''}`}>{status}</span><h3>{event.nome}</h3><p>{event.descricao || 'Sem descrição'}</p></button><div className="event-card-actions">{!event.publicado && !event.encerrado && <button className="action-icon publish-button" type="button" onClick={() => onPublish(event.id)} disabled={publishing} aria-label="Publicar evento" title="Publicar evento"><Send size={16} strokeWidth={2.2} aria-hidden="true" /></button>}{event.publicado && !event.encerrado && <button className="action-icon close-button" type="button" onClick={() => onClose(event.id)} disabled={closing} aria-label="Encerrar evento" title="Encerrar evento"><CircleStop size={16} strokeWidth={2.2} aria-hidden="true" /></button>}</div></article>
}
