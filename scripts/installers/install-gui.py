#!/usr/bin/env python3
"""
TormentNexus GUI Installer — Windows & macOS
Simple visual wizard that installs support for 38+ AI coding clients.
"""

import tkinter as tk
from tkinter import ttk, messagebox
import subprocess
import sys
import os
import threading

HOME = os.path.expanduser("~")
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
INSTALLER_PY = os.path.join(SCRIPT_DIR, "install-client-support.py")

WIDTH, HEIGHT = 700, 500
CLIENTS = [
    "Claude Code",
    "Gemini CLI",
    "Codex CLI",
    "Grok Build",
    "Antigravity",
    "Aider",
    "OpenCode",
    "OpenClaw",
    "Goose",
    "iFlow",
    "Roo Code",
    "Cline",
    "Cursor",
    "Windsurf",
    "Zed",
    "Trae",
    "Continue.dev",
    "Factory Droid",
    "OpenHands",
    "Kiro",
    "CodeWhale",
    "Omnigent",
    "Citadel",
    "Agent-Fusion",
    "Herdr",
    "Claude Squad",
    "Qwen Code",
    "Pi Coding Agent",
    "Kimi Code",
    "CLIProxyAPI",
    "VS Code",
    "JetBrains",
    "Hermes",
]


class InstallerApp:
    def __init__(self, root):
        self.root = root
        root.title("TormentNexus Installer")
        root.geometry(f"{WIDTH}x{HEIGHT}")
        root.resizable(False, False)
        try:
            root.iconbitmap(os.path.join(SCRIPT_DIR, "..", "go", "tormentnexus.exe"))
        except:
            pass

        # Colors
        self.BG = "#0a0a1a"
        self.CARD = "#151530"
        self.PRIMARY = "#667eea"
        self.TEXT = "#e0e0ff"
        self.SUBTEXT = "#8888bb"
        self.GREEN = "#00cc66"

        root.configure(bg=self.BG)

        self.build_ui()

    def build_ui(self):
        # Header
        header = tk.Frame(self.root, bg=self.BG)
        header.pack(fill="x", pady=(30, 10))

        title = tk.Label(
            header,
            text="TormentNexus",
            font=("Segoe UI", 28, "bold"),
            fg=self.PRIMARY,
            bg=self.BG,
        )
        title.pack()

        subtitle = tk.Label(
            header,
            text="Universal AI Coding Client Support Installer",
            font=("Segoe UI", 11),
            fg=self.SUBTEXT,
            bg=self.BG,
        )
        subtitle.pack(pady=(0, 5))

        tk.Label(
            header,
            text=f"Detected {len(CLIENTS)} supported AI coding agents",
            font=("Segoe UI", 9),
            fg="#6666aa",
            bg=self.BG,
        ).pack()

        # Client grid
        grid_frame = tk.Frame(self.root, bg=self.BG)
        grid_frame.pack(fill="both", expand=True, padx=40, pady=10)

        canvas = tk.Canvas(grid_frame, bg=self.BG, highlightthickness=0)
        scrollbar = tk.Scrollbar(grid_frame, orient="vertical", command=canvas.yview)
        self.client_frame = tk.Frame(canvas, bg=self.BG)

        canvas.create_window(
            (0, 0), window=self.client_frame, anchor="nw", tags="inner"
        )
        canvas.configure(yscrollcommand=scrollbar.set)

        # Grid layout: 3 columns
        for i, client in enumerate(CLIENTS):
            col = i % 3
            row = i // 3
            self._client_card(client, row, col)

        # Progress bar
        self.progress = ttk.Progressbar(self.root, mode="indeterminate", length=500)
        self.progress.pack(pady=(0, 5))

        # Status
        self.status_var = tk.StringVar(value="Ready to install")
        tk.Label(
            self.root,
            textvariable=self.status_var,
            font=("Segoe UI", 9),
            fg=self.SUBTEXT,
            bg=self.BG,
        ).pack()

        # Buttons
        btn_frame = tk.Frame(self.root, bg=self.BG)
        btn_frame.pack(pady=(10, 30))

        self.install_btn = tk.Button(
            btn_frame,
            text="Install Support for All Clients",
            font=("Segoe UI", 12, "bold"),
            bg=self.PRIMARY,
            fg="white",
            activebackground="#5577d9",
            activeforeground="white",
            relief="flat",
            padx=30,
            pady=10,
            cursor="hand2",
            command=self.install,
        )
        self.install_btn.pack(side="left", padx=5)

        tk.Button(
            btn_frame,
            text="Cancel",
            font=("Segoe UI", 11),
            bg=self.CARD,
            fg=self.SUBTEXT,
            activebackground="#202040",
            activeforeground=self.TEXT,
            relief="flat",
            padx=20,
            pady=10,
            cursor="hand2",
            command=self.root.destroy,
        ).pack(side="left", padx=5)

    def _client_card(self, name, row, col):
        card = tk.Frame(
            self.client_frame,
            bg=self.CARD,
            highlightbackground="#222255",
            highlightthickness=1,
            padx=10,
            pady=8,
        )
        card.grid(row=row, column=col, padx=4, pady=3, sticky="ew")

        tk.Label(
            card,
            text="\u2714",  # checkmark
            font=("Segoe UI", 10),
            fg=self.GREEN,
            bg=self.CARD,
        ).pack(side="left", padx=(0, 6))

        tk.Label(
            card,
            text=name,
            font=("Segoe UI", 10),
            fg=self.TEXT,
            bg=self.CARD,
        ).pack(side="left")

    def install(self):
        self.install_btn.configure(state="disabled", text="Installing...")
        self.progress.start(10)
        self.status_var.set("Installing support for all clients...")

        def run():
            try:
                result = subprocess.run(
                    [sys.executable, INSTALLER_PY],
                    capture_output=True,
                    text=True,
                    timeout=30,
                )
                self.root.after(0, lambda: self._done(True, result.stdout))
            except Exception as ex:
                self.root.after(0, lambda e=ex: self._done(False, str(e)))

        threading.Thread(target=run, daemon=True).start()

    def _done(self, success, output):
        self.progress.stop()
        self.progress.pack_forget()
        if success:
            self.install_btn.configure(
                text="Installation Complete!", bg=self.GREEN, state="normal"
            )
            self.status_var.set("38 clients installed successfully")
            messagebox.showinfo(
                "Success",
                "TormentNexus support installed for ALL clients!\n\n"
                "Start the kernel: tormentnexus serve\n"
                "Open dashboard: http://localhost:7779",
            )
        else:
            self.install_btn.configure(text="Try Again", bg="#ee4466", state="normal")
            self.status_var.set(f"Error: {output}")
            messagebox.showerror("Error", f"Installation failed: {output}")


def main():
    root = tk.Tk()
    app = InstallerApp(root)
    root.mainloop()


if __name__ == "__main__":
    main()
