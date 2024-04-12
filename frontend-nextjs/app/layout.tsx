import type { Metadata } from "next";
import "react-toastify/dist/ReactToastify.css";
import StoreProvider from "./components/StoreProvider";
import ToastLayout from "./components/ToastLayout";
import "./main.css";

export const metadata: Metadata = {
  title: "Fintrax - Financial App",
  description: "Fintrax - Financial App",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <StoreProvider>
        <body>
          {children}
          <ToastLayout />
        </body>
      </StoreProvider>
    </html>
  );
}
