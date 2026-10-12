import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { getOrder } from '../services/api'
import { useSession } from '../session/useSession'
import type { Order } from '../types'

export function CongressistaPurchaseSuccessPage() {
  const { session } = useSession()
  const { id } = useParams()
  const navigate = useNavigate()
  const [order, setOrder] = useState<Order | null>(null)
  useEffect(() => { if (session && id) void getOrder(session, Number(id)).then(setOrder) }, [session, id])
  return <div className="app-shell congressista-shell"><main className="main-content"><section className="checkout-view"><span className="kicker">Compra concluída</span><h1>Pagamento confirmado</h1>{order ? <><p>Pedido #{order.id} confirmado com sucesso. Seus ingressos estão disponíveis em Meus Eventos.</p><div className="cart-total"><span>Total</span><strong>{order.total.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })}</strong></div></> : <p>Carregando confirmação...</p>}  <button className="primary-button" onClick={() => navigate('/meus-eventos')}>Ver meus eventos</button></section></main></div>
}
