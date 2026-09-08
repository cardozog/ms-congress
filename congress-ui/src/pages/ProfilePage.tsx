import { useEffect, useState } from 'react'
import { getProfile } from '../services/api'
import { useSession } from '../session/useSession'
import type { Profile } from '../types'

export function ProfilePage() {
  const { session } = useSession(); const [profile, setProfile] = useState<Profile | null>(null)
  useEffect(() => { if (session) void getProfile(session).then(setProfile) }, [session])
  return <section className="view-content"><span className="kicker">Conta</span><p className="section-description">Dados usados para organizar seus eventos.</p><div className="profile-card"><Row label="E-mail" value={profile?.email ?? session?.email ?? 'Carregando...'} /><Row label="Razão social" value={profile?.pessoaJuridica?.razaoSocial ?? 'Não informado'} /><Row label="Nome fantasia" value={profile?.pessoaJuridica?.nomeFantasia ?? 'Não informado'} /><Row label="ID do organizador" value={String(profile?.pessoaJuridica?.id ?? 'Não informado')} /></div></section>
}
function Row({ label, value }: { label: string; value: string }) { return <div className="profile-row"><span>{label}</span><strong>{value}</strong></div> }
