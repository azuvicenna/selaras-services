#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.5cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123 | Email: rs\@selaras.hospital.id]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Dokumen Rujukan Medis Resmi - Terverifikasi Elektronik]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 10pt)

#align(center)[
  #text(12pt, weight: "bold")[SURAT RUJUKAN PASIEN] \
  #text(9pt)[Nomor: {{DOC_NUMBER}}]
]

#v(1em)

Kepada Yth. \
*TS Dokter Spesialis {{SPESIALISASI_TUJUAN}}* \
di *{{RS_TUJUAN}}*

#v(0.5em)
Dengan hormat, \
Bersama ini kami merujuk pasien dengan data identitas sebagai berikut:

#table(
  columns: (140pt, 10pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}],
  [Nomor Rekam Medis], [:], [{{PATIENT_ID}}],
  [Jenis Kelamin / Tgl Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}],
)

#v(0.5em)
#table(
  columns: (1fr),
  stroke: 0.5pt,
  [ *Anamnesa & Pemeriksaan Fisik:* \ {{ANAMNESA}} ],
  [ *Diagnosa Sementara:* \ {{DIAGNOSIS}} ],
  [ *Tindakan & Terapi yang Telah Diberikan:* \ {{TERAPI_SEBELUMNYA}} ],
  [ *Alasan Rujukan / Permintaan:* \ {{ALASAN_RUJUKAN}} ],
)

#v(1em)
Mohon evaluasi dan penatalaksanaan lebih lanjut. Atas kerja sama dan bantuan sejawat, kami sampaikan terima kasih.

#v(2em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI KEASLIAN ELEKTRONIK] \
        #text(6.5pt)[Pindai kode QR atau akses URL di bawah untuk verifikasi keabsahan dokumen:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Pengirim, \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
