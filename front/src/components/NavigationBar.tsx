import React, { useState } from "react";
import { User } from "../types/user";
import { NavLink } from "react-router";

interface NavigationBarProps {
  user: User | null;
  logout: () => void;
}

const NavigationBar: React.FC<NavigationBarProps> = ({ user, logout }) => {
  const [menuOpen, setMenuOpen] = useState(false);

  const navLinks = [
    { to: "/products", label: "Products", show: true },
    { to: "/users/" + (user?.id || ""), label: "Profile", show: !!user },
    { to: "/users/" + (user?.id || "") + "/cart", label: "Cart", show: !!user },
    { to: "/users/" + (user?.id || "") + "/orders", label: "Orders", show: !!user },
  ];

  return (
    <nav className="bg-black/95 backdrop-blur-sm border-b border-purple-800/50 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo/Brand */}
          <NavLink to="/" className="flex items-center select-none">
            <span className="text-2xl font-bold text-white">
              Bachelor<span className="text-purple-400">Shop</span>
            </span>
          </NavLink>

          {/* Desktop Navigation */}
          <div className="hidden md:flex md:items-center md:space-x-6">
            {navLinks.filter(l => l.show).map(link => (
              <NavLink
                key={link.to}
                to={link.to}
                className={({ isActive }) =>
                  `px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 w-24 text-center
                  ${isActive
                    ? "bg-purple-900/50 text-purple-200"
                    : "text-gray-300 hover:bg-purple-800/30 hover:text-purple-200"}`
                }
              >
                {link.label}
              </NavLink>
            ))}

            {!user ? (
              <div className="flex items-center space-x-4 ml-6">
                <NavLink
                  to="/login"
                  className="px-4 py-2 rounded-md text-sm font-medium bg-purple-600 hover:bg-purple-700 text-white transition-colors duration-200 w-24 text-center"
                >
                  Login
                </NavLink>
                <NavLink
                  to="/register"
                  className="px-4 py-2 rounded-md text-sm font-medium border border-purple-500/50 text-purple-300 hover:bg-purple-900/30 transition-colors duration-200 w-24 text-center"
                >
                  Register
                </NavLink>
              </div>
            ) : (
              <button
                onClick={logout}
                className="ml-6 px-4 py-2 rounded-md text-sm font-medium bg-zinc-800/80 text-purple-200 hover:bg-purple-900/30 transition-colors duration-200 border border-purple-700/50 w-24 text-center"
              >
                Logout
              </button>
            )}
          </div>

          {/* Mobile menu button */}
          <button
            onClick={() => setMenuOpen(!menuOpen)}
            className="md:hidden p-2 rounded-md text-purple-300 hover:bg-purple-900/30 focus:outline-none"
            aria-label="Toggle menu"
          >
            <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
        </div>
      </div>

      {/* Mobile menu */}
      {menuOpen && (
        <div className="md:hidden bg-black/95 border-t border-purple-800/50">
          <div className="px-2 pt-2 pb-3 space-y-1">
            {navLinks.filter(l => l.show).map(link => (
              <NavLink
                key={link.to}
                to={link.to}
                className={({ isActive }) =>
                  `block px-3 py-2 rounded-md text-base font-medium transition-colors duration-200
                  ${isActive
                    ? "bg-purple-900/50 text-purple-200"
                    : "text-gray-300 hover:bg-purple-800/30 hover:text-purple-200"}`
                }
                onClick={() => setMenuOpen(false)}
              >
                {link.label}
              </NavLink>
            ))}

            {!user ? (
              <div className="mt-4 space-y-2">
                <NavLink
                  to="/login"
                  className="block w-full px-3 py-2 rounded-md text-base font-medium bg-purple-600 hover:bg-purple-700 text-white text-center transition-colors duration-200"
                  onClick={() => setMenuOpen(false)}
                >
                  Login
                </NavLink>
                <NavLink
                  to="/register"
                  className="block w-full px-3 py-2 rounded-md text-base font-medium border border-purple-500/50 text-purple-300 hover:bg-purple-900/30 text-center transition-colors duration-200"
                  onClick={() => setMenuOpen(false)}
                >
                  Register
                </NavLink>
              </div>
            ) : (
              <button
                onClick={() => { setMenuOpen(false); logout(); }}
                className="w-full mt-4 px-3 py-2 rounded-md text-base font-medium bg-zinc-800/80 text-purple-200 hover:bg-purple-900/30 transition-colors duration-200 border border-purple-700/50"
              >
                Logout
              </button>
            )}
          </div>
        </div>
      )}
    </nav>
  );
};

export default NavigationBar; 