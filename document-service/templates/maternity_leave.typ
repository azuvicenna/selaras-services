#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.5cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(10pt, weight: "bold")[POLIKLINIK KEBIDANAN & KANDUNGAN] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Surat Keterangan Cuti Bersalin Resmi - RS Selaras Medika]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 10pt)

#align(center)[
  #text(12pt, weight: "bold")[SURAT KETERANGAN CUTI HAMIL / BERSALIN] \
  #text(9pt)[Nomor: {{DOC_NUMBER}}]
]

#v(1em)

Yang bertanda tangan di bawah ini Dokter Spesialis Kebidanan dan Kandungan menerangkan bahwa:

#table(
  columns: (140pt, 10pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}],
  [Nomor Rekam Medis], [:], [{{PATIENT_ID}}],
  [Tgl Lahir / Umur], [:], [{{PATIENT_DOB}} / {{UMUR}} Tahun],
  [Pekerjaan / Instansi], [:], [{{INSTANSI}}],
)

#v(0.5em)
Berdasarkan hasil pemeriksaan obstetri terkini pada tanggal {{ISSUED_DATE}}:
#table(
  columns: (150pt, 10pt, 1fr),
  stroke: 0.5pt,
  [Usia Kehamilan (Gestasi)], [:], [{{USIA_KEHAMILAN}} Minggu],
  [Hari Perkiraan Lahir (HPL)], [:], [{{HPL}}],
  [Gravida / Para / Abortus], [:], [G{{G}}P{{P}}A{{A}}],
  [Diagnosa Klinis], [:], [{{DIAGNOSIS}}],
)

#v(1em)
Mengingat kondisi kehamilan tersebut di atas, pasien disarankan untuk mengambil *Cuti Hamil dan Bersalin* selama:

#align(center)[
  #rect(stroke: 1pt, inset: 10pt)[
    #text(12pt, weight: "bold")[3 (Tiga) Bulan] \
    #text(9pt)[Terhitung mulai tanggal {{TGL_MULAI_CUTI}} s/d {{TGL_SELESAI_CUTI}}]
  ]
]

#v(1em)
Demikian surat keterangan ini diberikan untuk digunakan sebagai persyaratan pengajuan cuti bersalin pada instansi terkait.

#v(2em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI KEASLIAN ELEKTRONIK] \
        #text(6.5pt)[Pindai kode QR untuk validasi keaslian surat:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Spesialis Kebidanan (Sp.OG), \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
