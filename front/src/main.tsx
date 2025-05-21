import {
  createBrowserRouter,
  RouterProvider,
} from "react-router";

import ReactDOM from "react-dom/client";

import "./index.css";

// Import your actual components here
import MainLayout from "./components/MainLayout";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Users from "./pages/Users";
import UserDetails from "./pages/UserDetails";
import UserActivities from "./pages/UserActivities";
import UserFavorites from "./pages/UserFavorites";
import UserCart from "./pages/UserCart";
import UserOrders from "./pages/UserOrders";
import Products from "./pages/Products";
import ProductDetails from "./pages/ProductDetails";
import ProductsByCategory from "./pages/ProductsByCategory";
import CartDetails from "./pages/CartDetails";
import CartItemEdit from "./pages/CartItemEdit";
import Orders from "./pages/Orders";
import OrderDetails from "./pages/OrderDetails";
import AnalyticsActivities from "./pages/AnalyticsActivities";
import AnalyticsRequests from "./pages/AnalyticsRequests";
import AnalyticsPopularProducts from "./pages/AnalyticsPopularProducts";
import AnalyticsActiveUsers from "./pages/AnalyticsActiveUsers";
import NotFound from "./pages/NotFound";
import RequireAuth from "./components/RequireAuth";

const router = createBrowserRouter([
  {
    path: "/",
    element: <MainLayout />, 
    children: [
      { index: true, element: <Home /> },
      { path: "login", element: <Login /> },
      { path: "register", element: <Register /> },
      {
        path: "users",
        element: <RequireAuth><Users /></RequireAuth>,
        children: [
          { index: true, element: <RequireAuth><Users /></RequireAuth> },
          { path: ":id", element: <RequireAuth><UserDetails /></RequireAuth> },
          { path: ":id/activities", element: <RequireAuth><UserActivities /></RequireAuth> },
          { path: ":id/favorites", element: <RequireAuth><UserFavorites /></RequireAuth> },
          { path: ":id/cart", element: <RequireAuth><UserCart /></RequireAuth> },
          { path: ":id/orders", element: <RequireAuth><UserOrders /></RequireAuth> },
        ],
      },
      {
        path: "products",
        children: [
          { index: true, element: <Products /> },
          { path: ":id", element: <ProductDetails /> },
          { path: "category/:category", element: <ProductsByCategory /> },
        ],
      },
      {
        path: "cart",
        children: [
          { path: ":id", element: <RequireAuth><CartDetails /></RequireAuth> },
          { path: ":id/items/:itemId", element: <RequireAuth><CartItemEdit /></RequireAuth> },
        ],
      },
      {
        path: "orders",
        children: [
          { index: true, element: <RequireAuth><Orders /></RequireAuth> },
          { path: ":id", element: <RequireAuth><OrderDetails /></RequireAuth> },
        ],
      },
      {
        path: "analytics",
        children: [
          { path: "activities", element: <RequireAuth><AnalyticsActivities /></RequireAuth> },
          { path: "requests", element: <RequireAuth><AnalyticsRequests /></RequireAuth> },
          { path: "products/popular", element: <RequireAuth><AnalyticsPopularProducts /></RequireAuth> },
          { path: "users/active", element: <RequireAuth><AnalyticsActiveUsers /></RequireAuth> },
        ],
      },
      { path: "*", element: <NotFound /> }, // 404 fallback
    ],
  },
]);

const root = document.getElementById("root");
if (!root) throw new Error("Root element not found");

ReactDOM.createRoot(root).render(
  <RouterProvider router={router} />
);