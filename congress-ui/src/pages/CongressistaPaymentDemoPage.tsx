import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { getOrder, simulatePayment, startPayment } from '../services/api'
import { useSession } from '../session/useSession'
import type { Order } from '../types'
import { errorMessage } from '../utils'

const currency = (value: number) => value.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })

export function CongressistaPaymentDemoPage() {
  const { session } = useSession()
  const { id } = useParams()
  const navigate = useNavigate()
  const [order, setOrder] = useState<Order | null>(null)
  const [feedback, setFeedback] = useState('')
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    if (!session || !id) return
    void getOrder(session, Number(id)).then(setOrder).catch(error => setFeedback(errorMessage(error, 'Não foi possível carregar o pedido.')))
  }, [session, id])

  async function begin() {
    if (!session || !order) return
    setBusy(true)
    try { setOrder(await startPayment(session, order.id, crypto.randomUUID())) }
    catch (error) { setFeedback(errorMessage(error, 'Não foi possível iniciar o pagamento.')) }
    finally { setBusy(false) }
  }

  async function simulate(approved: boolean) {
    if (!session || !order) return
    setBusy(true)
    try {
      const updated = await simulatePayment(session, order.id, approved)
      if (approved && updated.status === 'PAID') navigate(`/congressista/compra-concluida/${order.id}`)
      else setOrder(updated)
    } catch (error) { setFeedback(errorMessage(error, 'Não foi possível simular o pagamento.')) }
    finally { setBusy(false) }
  }

  if (!order) return <div className="app-shell congressista-shell"><main className="main-content"><div className="empty-state"><p>{feedback || 'Carregando pagamento...'}</p></div></main></div>
  return <div className="app-shell congressista-shell"><main className="main-content"><section className="checkout-view"><span className="kicker">Demonstração de pagamento</span><h1>Pedido #{order.id}</h1><p>Este pagamento é simulado exclusivamente para demonstração. Nenhum dado real de cartão é coletado.</p><div className="cart-total"><span>Valor total</span><strong>{currency(order.total)}</strong></div>{feedback && <p className="feedback">{feedback}</p>}{!order.pagamento?.referenciaExterna ? <button className="primary-button" disabled={busy} onClick={() => void begin()}>{busy ? 'Iniciando...' : 'Preparar pagamento simulado'}</button> : <div className="payment-actions"><button className="primary-button" disabled={busy} onClick={() => void simulate(true)}>Simular pagamento aprovado</button><button className="secondary-button" disabled={busy} onClick={() => void simulate(false)}>Simular pagamento recusado</button></div>}<button className="secondary-button" onClick={() => navigate('/congressista/carrinho')}>Voltar ao carrinho</button></section></main></div>
}
