import { startTransition, useEffect, useState, type FormEvent } from 'react'
import type { EventDetails, EventPayload, TicketPayload, TicketType } from '../types'
import { fromInputDate, toInputDate } from '../utils'

interface EventFormProps { event?: EventDetails | null; ticketTypes: TicketType[]; organizerId: number; loading: boolean; feedback: string; onSubmit: (payload: EventPayload) => Promise<void>; onCancel: () => void }
type FormValues = Omit<EventPayload, 'data_inicio' | 'data_fim' | 'organizador_id' | 'ingressos'> & { data_inicio: string; data_fim: string; ingressos: TicketPayload[] }
const empty: FormValues = { nome: '', descricao: '', data_inicio: '', data_fim: '', logradouro: '', cidade: '', estado: '', cep: '', ingressos: [{ tipo_ingresso_id: 0, preco: 0, quantidade: 1 }] }

function eventValues(event?: EventDetails | null): FormValues {
  if (!event) return empty
  return { nome: event.nome, descricao: event.descricao, data_inicio: toInputDate(event.dataInicio), data_fim: toInputDate(event.dataFim), logradouro: event.logradouro, cidade: event.cidade, estado: event.estado, cep: event.cep, ingressos: event.ingressos?.map((ticket) => ({ tipo_ingresso_id: ticket.tipoIngressoID, preco: ticket.preco, quantidade: ticket.quantidade })) ?? empty.ingressos }
}

export function EventForm({ event, ticketTypes, organizerId, loading, feedback, onSubmit, onCancel }: EventFormProps) {
  const [values, setValues] = useState<FormValues>(() => eventValues(event))
  useEffect(() => { startTransition(() => setValues(eventValues(event))) }, [event])
  function update(name: keyof Omit<FormValues, 'ingressos'>, value: string) { setValues((current) => ({ ...current, [name]: value })) }
  function updateTicket(index: number, name: keyof TicketPayload, value: string) { setValues((current) => ({ ...current, ingressos: current.ingressos.map((ticket, ticketIndex) => ticketIndex === index ? { ...ticket, [name]: Number(value) } : ticket) })) }
  function addTicket() { setValues((current) => ({ ...current, ingressos: [...current.ingressos, { tipo_ingresso_id: 0, preco: 0, quantidade: 1 }] })) }
  function removeTicket(index: number) { setValues((current) => ({ ...current, ingressos: current.ingressos.filter((_, ticketIndex) => ticketIndex !== index) })) }
  async function submit(formEvent: FormEvent) {
    formEvent.preventDefault()
    if (new Date(values.data_fim) <= new Date(values.data_inicio) || values.ingressos.some((ticket) => ticket.tipo_ingresso_id === 0)) return
    await onSubmit({ ...values, data_inicio: fromInputDate(values.data_inicio), data_fim: fromInputDate(values.data_fim), estado: values.estado.toUpperCase(), organizador_id: organizerId })
  }
  return <div className="form-view"><button className="back-button" type="button" onClick={onCancel}>← Voltar para eventos</button><div className="form-view-heading"><span className="kicker">{event ? 'Atualização' : 'Novo evento'}</span><h2>{event ? 'Editar evento' : 'Criar evento'}</h2></div>
    <form className="event-form" onSubmit={submit}><div className="form-grid"><Field label="Nome do evento" name="nome" value={values.nome} onChange={update} wide required /><Field label="Descrição" name="descricao" value={values.descricao} onChange={update} wide textarea required /><Field label="Início" name="data_inicio" value={values.data_inicio} onChange={update} type="datetime-local" required /><Field label="Fim" name="data_fim" value={values.data_fim} onChange={update} type="datetime-local" required /><Field label="Logradouro" name="logradouro" value={values.logradouro} onChange={update} wide required /><Field label="Cidade" name="cidade" value={values.cidade} onChange={update} required /><Field label="Estado" name="estado" value={values.estado} onChange={update} required /><Field label="CEP" name="cep" value={values.cep} onChange={update} required /></div>
      <div className="tickets-section"><div className="tickets-heading"><div><span className="kicker">Ingressos</span><p>Selecione os tipos cadastrados e defina preço e quantidade.</p></div><button className="secondary-button" type="button" onClick={addTicket}>+ Adicionar tipo</button></div>{values.ingressos.map((ticket, index) => <div className="ticket-row" key={`${index}-${ticket.tipo_ingresso_id}`}><div className="field"><label htmlFor={`ticket-type-${index}`}>Tipo</label><select id={`ticket-type-${index}`} value={ticket.tipo_ingresso_id} onChange={(formEvent) => updateTicket(index, 'tipo_ingresso_id', formEvent.target.value)} required><option value={0}>Selecione</option>{ticketTypes.map((type) => <option key={type.id} value={type.id}>{type.nome}</option>)}</select></div><Field label="Preço" name="preco" value={String(ticket.preco)} onChange={(_, value) => updateTicket(index, 'preco', value)} type="number" required /><Field label="Quantidade" name="quantidade" value={String(ticket.quantidade)} onChange={(_, value) => updateTicket(index, 'quantidade', value)} type="number" required />{values.ingressos.length > 1 && <button className="remove-ticket" type="button" onClick={() => removeTicket(index)} aria-label="Remover tipo de ingresso">×</button>}</div>)}</div>
      {feedback && <p className="feedback">{feedback}</p>}<div className="form-actions"><button className="secondary-button" type="button" onClick={onCancel}>Cancelar</button><button className="primary-button" disabled={loading}>{loading ? 'Salvando...' : event ? 'Salvar alterações' : 'Publicar evento'} <b>→</b></button></div>
    </form></div>
}

interface FieldProps { label: string; name: string; value: string; onChange: (name: never, value: string) => void; wide?: boolean; textarea?: boolean; type?: string; required?: boolean }
function Field({ label, name, value, onChange, wide, textarea, type = 'text', required }: FieldProps) { return <div className={`field${wide ? ' field-wide' : ''}`}><label htmlFor={name}>{label}</label>{textarea ? <textarea id={name} value={value} onChange={(event) => onChange(name as never, event.target.value)} rows={4} required={required} /> : <input id={name} type={type} min={type === 'number' ? 0 : undefined} step={name === 'preco' ? '0.01' : undefined} value={value} onChange={(event) => onChange(name as never, event.target.value)} required={required} />}</div> }
