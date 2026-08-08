import { createBrowserRouter, Navigate } from 'react-router-dom';
import { ProfilesPage } from '@/pages/profiles/ProfilesPage';
import { GeneratingPage } from '@/pages/generating/GeneratingPage';
import { RecapPage } from '@/pages/recap/RecapPage';

/**
 * Три шага сценария — три маршрута. Разделены намеренно: на генерацию и на
 * готовый recap можно дать прямую ссылку, а «назад» в браузере не ломает историю.
 */
export const router = createBrowserRouter([
  { path: '/', element: <ProfilesPage /> },
  { path: '/generate/:profileId', element: <GeneratingPage /> },
  { path: '/recap/:profileId', element: <RecapPage /> },
  { path: '*', element: <Navigate to="/" replace /> },
]);
