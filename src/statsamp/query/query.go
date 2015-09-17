package query

import (
  "bytes"
  "net"
  "time"
  "encoding/binary"
)

type ServerInfo struct {
  Password    bool
  Players     int
  MaxPlayers  int
  Hostname    string
  Gamemode    string
  MapName     string
}

func GetServerInfo(serverAddr string, timeout int) (ServerInfo, error) {
  addr, err := net.ResolveUDPAddr("udp", serverAddr)
  if err != nil {
    return ServerInfo{}, err
  }
  con, err := net.DialUDP("udp", nil, addr)
  if err != nil {
    return ServerInfo{}, err
  }
  // Writing to socket
  con.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
  query := encodeQuery(addr.IP, addr.Port, "i")
  if _, err := con.Write(query); err != nil {
    return ServerInfo{}, err
  }
  // Reading from socket
  con.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
  response := make([]byte, 512)
  if _, err := con.Read(response); err != nil {
    return ServerInfo{}, err
  }
  return decodeServerInfo(response)
}

func encodeQuery(ip net.IP, port int, queryType string) []byte {
  buf := new(bytes.Buffer)

  buf.WriteString("SAMP")
  buf.Write(ip[len(ip)-4:]) // IP in 4 bytes
  buf.WriteByte(byte(port & 0xFF)) // First byte of port
  buf.WriteByte(byte(port >> 8 & 0xFF)) // Second byte of port
  buf.WriteString(queryType) // Query type (i, r, c, d, x, p)

  return buf.Bytes()
}

func decodeServerInfo(response []byte) (ServerInfo, error) {
  var info ServerInfo
  buf := bytes.NewBuffer(response[11:])
  // Password
  password, err := buf.ReadByte()
  if err != nil {
    return info, err
  }
  info.Password = (password != 0x0)
  // Players
  players, err := readInt16(buf)
  if err != nil {
    return info, err
  }
  info.Players = int(players)
  // MaxPlayers
  maxPlayers, err := readInt16(buf)
  if err != nil {
    return info, err
  }
  info.MaxPlayers = int(maxPlayers)
  // Hostname
  info.Hostname, err = readString(buf)
  if err != nil {
    return info, err
  }
  info.Hostname, err = transformToUTF8(info.Hostname)
  if err != nil {
    return info, err
  }
  // Gamemode
  info.Gamemode, err = readString(buf)
  if err != nil {
    return info, err
  }
  info.Gamemode, err = transformToUTF8(info.Gamemode)
  if err != nil {
    return info, err
  }
  // Mapname
  info.MapName, err = readString(buf)
  if err != nil {
    return info, err
  }
  info.MapName, err = transformToUTF8(info.MapName)
  if err != nil {
    return info, err
  }
  return info, nil
}

func readInt16(buf *bytes.Buffer) (int16, error) {
  var result int16
  err := binary.Read(buf, binary.LittleEndian, &result)
  return result, err
}

func readString(buf *bytes.Buffer) (string, error) {
  var length int32
  if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
    return "", err
  }
  result := make([]byte, length)
  if err := binary.Read(buf, binary.LittleEndian, &result); err != nil {
    return "", err
  }
  return string(result), nil
}
