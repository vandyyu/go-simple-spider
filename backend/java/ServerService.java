package server;
import java.io.InputStreamReader;
import java.io.BufferedReader;
import java.io.OutputStreamWriter;
import java.io.BufferedWriter;
import java.net.Socket;
import java.io.ByteArrayInputStream;
import java.io.OutputStream;

class ServerService{
	private Socket socket;
	public ServerService(Socket socket){
		this.socket = socket;
	}
	public void Run(SeleniumPool pool){
		BufferedReader reader = null;
		try{
			reader = new BufferedReader(new InputStreamReader(this.socket.getInputStream()));
			String line = reader.readLine();
			if(line != null){
				this.readWrite(pool, line);
			}else{
				this.readWrite(pool, "");
			}
		}catch(Exception e){
			e.printStackTrace();
		}finally{
			try{
				if (reader != null){
					reader.close();
				}
				if(this.socket != null){
					this.socket.close();
					this.socket = null;
				}
			}catch(Exception e){
				e.printStackTrace();
			}
		}
	}
	private void readWrite(SeleniumPool pool, String line){
		OutputStream writer = null;
		try{
			writer = this.socket.getOutputStream();
			String html = pool.getFreeExecutor().getHTMLData(line);
			if (html == null){
				writer.write("".getBytes());
			}else{
				ByteArrayInputStream bis = new ByteArrayInputStream(html.getBytes());
				int len = 0;
				byte[] buf = new byte[1024];
				while((len = bis.read(buf, 0, buf.length)) > 0){
					writer.write(buf, 0, len);
				}

				// do not write data by once. because [java.net.SocketException: Broken pipe] will occur when too much data to transfer.
				// writer.write(html);
			}
		}catch(Exception e){
			e.printStackTrace();
		}finally{
			
			try{
				if(writer != null){
					writer.flush();
					writer.close();
				}
			}catch(Exception e){
				e.printStackTrace();
			}
			
		}
	}
}
