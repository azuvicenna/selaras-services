#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Formulir Persetujuan Tindakan Kedokteran (Informed Consent)]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9pt)

#align(center)[
  #text(12pt, weight: "bold")[SURAT PERSETUJUAN TINDAKAN MEDIS (INFORMED CONSENT)] \
  #text(8.5pt)[Nomor: {{DOC_NUMBER}}]
]

#v(0.5em)

*I. PEMBERIAN INFORMASI & EDUKASI OLEH DOKTER:*
#table(
  columns: (120pt, 8pt, 1fr),
  stroke: 0.5pt,
  [Dokter Pelaksana Tindakan], [:], [{{DOCTOR_NAME}}],
  [Pemberi Informasi / Edukasi], [:], [{{PEMBERI_INFORMASI}}],
  [Penerima Informasi / Edukasi], [:], [{{PENERIMA_INFORMASI}} (Hubungan: {{HUBUNGAN}})],
  [Diagnosis Kerja / Banding], [:], [{{DIAGNOSIS}}],
  [Rencana Tindakan Medis], [:], [{{TINDAKAN_MEDIS}}],
  [Indikasi & Tata Cara Tindakan], [:], [{{TATA_CARA}}],
  [Tujuan Tindakan], [:], [{{TUJUAN}}],
  [Risiko & Komplikasi Tindakan], [:], [{{RISIKO}}],
  [Prognosis], [:], [{{PROGNOSIS}}],
  [Alternatif Tindakan & Risiko], [:], [{{ALTERNATIF}}],
)

#v(0.5em)
*II. PERNYATAAN PERSETUJUAN PASIEN / KELUARGA:*
Yang bertanda tangan di bawah ini:
#table(
  columns: (120pt, 8pt, 1fr),
  stroke: none,
  [Nama], [:], [{{PENERIMA_INFORMASI}}],
  [Jenis Kelamin / Tgl Lahir], [:], [{{GENDER_WALI}} / {{DOB_WALI}}],
  [Alamat / No. Telp], [:], [{{ALAMAT_WALI}}],
  [Hubungan dengan Pasien], [:], [{{HUBUNGAN}}],
)

Menyatakan dengan sesungguhnya telah *MEMAHAMI SEPENUHNYA* seluruh penjelasan medis dan dengan ini memberikan *PERSETUJUAN* untuk dilakukannya tindakan kedokteran tersebut di atas terhadap diri saya / pasien bernama: *{{PATIENT_NAME}}* (No. RM: *{{PATIENT_ID}}*).

#v(1em)

#grid(
  columns: (1fr, 1fr, 1fr),
  [
    #align(center)[
      Dokter Pelaksana, \
      #v(3em)
      #text(weight: "bold")[{{DOCTOR_NAME}}]
    ]
  ],
  [
    #align(center)[
      Saksi Pihak RS (Perawat), \
      #v(3em)
      #text(weight: "bold")[{{SAKSI_RS}}]
    ]
  ],
  [
    #align(center)[
      Yang Membuat Pernyataan, \
      #v(3em)
      #text(weight: "bold")[{{PENERIMA_INFORMASI}}]
    ]
  ]
)
