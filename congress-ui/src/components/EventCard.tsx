import type { EventSummary } from '../types'
import { formatDate } from '../utils'

export function EventCard({ event, onSelect }: { event: EventSummary; onSelect: (id: number) => void }) {
  return <article className="event-card"><button type="button" onClick={() => onSelect(event.id)}><span className="event-date">{formatDate(event.dataInicio)}</span><h3>{event.nome}</h3><p>{event.descricao || 'Sem descrição'}</p></button></article>
}
