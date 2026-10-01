             wordlist
                │
                ▼
        ┌────────────────┐
        │      ffuf      │
        └───────┬────────┘
                │
       replace FUZZ
                │
                ▼
        HTTP request
                │
                ▼
          Web server
                │
                ▼
          HTTP response
                │
                ▼
       filters/matchers
                │
                ▼
       interesting results