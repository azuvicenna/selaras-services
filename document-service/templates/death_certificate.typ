#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(10pt, weight: "bold")[INSTALASI KEDOKTERAN FORENSIK & MEDIKOLEGAL] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Surat Keterangan Kematian Resmi - RS Selaras Medika]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9.5pt)

#align(center)[
  #text(12pt, weight: "bold")[SURAT KETERANGAN KEMATIAN / VISUM] \
  #text(9pt)[Nomor: {{DOC_NUMBER}}]
]

#v(0.8em)

Yang bertanda tangan di bawah ini dokter Rumah Sakit Selaras Medika menerangkan bahwa:

#table(
  columns: (140pt, 10pt, 1fr),
  stroke: none,
  [Nama Jenazah / Almarhum], [:], [{{PATIENT_NAME}}],
  [Nomor Rekam Medis], [:], [{{PATIENT_ID}}],
  [Nomor Induk Kependudukan (NIK)], [:], [{{NIK}}],
  [Jenis Kelamin / Umur], [:], [{{PATIENT_GENDER}} / {{UMUR}} Tahun],
  [Alamat Domisili Terakhir], [:], [{{ALAMAT}}],
)

#v(0.5em)
Telah dinyatakan *MENINGGAL DUNIA* di Rumah Sakit Selaras Medika pada:
#table(
  columns: (140pt, 10pt, 1fr),
  stroke: 0.5pt,
  [Hari / Tanggal Kematian], [:], [{{HARI_KEMATIAN}}, {{TGL_KEMATIAN}}],
  [Waktu Kematian], [:], [Pukul {{JAM_KEMATIAN}} WIB],
  [Tempat Kematian], [:], [{{TEMPAT_KEMATIAN}} (IGD / ICU / Ruang Rawat)],
)

#v(0.5em)
*PENYEBAB KEMATIAN (Penyebab Medis / ICD-10):*
#table(
  columns: (160pt, 10pt, 1fr),
  stroke: 0.5pt,
  [I. a. Sebab Langsung (Immediate Cause)], [:], [{{SEBAB_LANGSUNG}}],
  [   b. Sebab Antara (Intermediate Cause)], [:], [{{SEBAB_ANTARA}}],
  [   c. Sebab Dasar (Underlying Cause)], [:], [{{SEBAB_DASAR}}],
  [II. Kondisi Lain yang Berkontribusi], [:], [{{KONDISI_KONTRIBUSI}}],
)

#v(0.5em)
*Catatan / Visum et Repertum:* \
{{VISUM_NOTES}}

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI KEASLIAN ELEKTRONIK] \
        #text(6.5pt)[Pindai kode QR untuk validasi keaslian surat kematian:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter yang Memeriksa / Forensik, \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
