import { createBrowserRouter, Navigate } from 'react-router-dom';
import { ProfilesPage } from '@/pages/profiles/ProfilesPage';
import { GeneratingPage } from '@/pages/generating/GeneratingPage';
import { RecapPage } from '@/pages/recap/RecapPage';
import { PublicSharePage } from '@/pages/share/PublicSharePage';

/**
 * Три шага сценария — три маршрута. Итоги адресуются по `recapId`, а не по
 * профилю: snapshot неизменяем, поэтому ссылка на него переживает перезагрузку
 * и не запускает генерацию заново.
 */
export const router = createBrowserRouter([
  { path: '/', element: <ProfilesPage /> },
  { path: '/generate/:profileId/:year', element: <GeneratingPage /> },
  { path: '/recap/:recapId', element: <RecapPage /> },
  // Публичный просмотр: отдельный маршрут, чтобы ссылка наружу вела не на
  // приватные итоги, а на безопасную проекцию из /share.
  { path: '/share/:recapId', element: <PublicSharePage /> },
  { path: '*', element: <Navigate to="/" replace /> },
]);
