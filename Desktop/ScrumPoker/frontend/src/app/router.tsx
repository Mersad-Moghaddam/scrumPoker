import { createBrowserRouter } from "react-router-dom";

import { HomePage } from "../features/home/home-page";
import { RoomPage } from "../features/room/room-page";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <HomePage />
  },
  {
    path: "/room/:code",
    element: <RoomPage />
  }
]);
