package server;
import java.net.Socket;
import java.net.ServerSocket;
import java.util.List;
import java.util.ArrayList;

class SpiderServer{
	private List<Thread> threads;
	
	private static int POOL_SIZE = 5;

	private ServerSocket serverSocket ;
	private boolean stop;

	public SpiderServer(int port){
		try{	
			this.serverSocket = new ServerSocket(port);
		}catch(Exception e){
			e.printStackTrace();
		}
		this.threads = new ArrayList<Thread>();
		this.stop = false;
	}
	public void listen(){
		SeleniumPool pool = new SeleniumPool(POOL_SIZE);
		pool.initFirefoxPool();
		while(!stop){
			try{
				Socket socket = this.serverSocket.accept();
				Runnable r = new Runnable(){
					public void run(){
						ServerService service = new ServerService(socket);
						service.Run(pool);
					}
				};
				Thread t = new Thread(r);
				this.threads.add(t);
				t.start();
			}catch(Exception e){
				e.printStackTrace();
			}
		}
		for(int i = 0;i < threads.size();i++){
			Thread t = threads.get(i);
			try{
				t.join();
			}catch(Exception e){
				e.printStackTrace();
			}
		}
		pool.destroy();
	}

	// TODO: stop server, change stop value.
}
